package epaxos

import (
	"fmt"

	"github.com/bbengfort/epaxos/pb"
)

// NewLog creates a new 2D log for epaxos.
func NewLog(config *Config) *Logs {
	logs := new(Logs)
	logs.logs = make(map[uint32]*replicaLog)

	for _, peer := range config.Peers {
		logs.logs[peer.PID] = &replicaLog{
			conflicts: make(map[string]uint64),
			instances: make([]*pb.Instance, 0),
		}
	}

	return logs
}

// Logs is a 2D array that maintains the state of all replicas by keeping track of
// Instances of operations that can be applied to the cluster state. This type exposes
// a number of helpful methods to interact with the log and sequence numbers.
//
// An instance in a log is identified by (pid, slot) where the PID is the id of the
// replica who is the quorum leader for the instance and the slot is the index in the
// log for that replica where the instance resides. Instances are further defined by
// sequence numbers and dependencies. Sequence numbers expose the order that the
// instance was applied to this particular log, and dependencies order all instnaces
// with respect to all internal logs.
type Logs struct {
	logs     map[uint32]*replicaLog // the 2D internal log slices managed by replica PID
	sequence uint64                 // the maximum sequence number seen by this log
}

// An internal type for the slice of Instances assigned to each replica.
type replicaLog struct {
	conflicts map[string]uint64 // cache of key to latest instance to optimize conflict detection
	instances []*pb.Instance    // the instances created by the replica
}

//===========================================================================
// Instances Management
//===========================================================================

// Create an instance based on the current log state and append it. This is a helper
// method for quickly adding a set of operations to the log in a consistent manner.
func (l *Logs) Create(replica uint32, ops []*pb.Operation) (inst *pb.Instance, err error) {
	// Increment the sequence number
	l.sequence++

	// Get the log of the specified replica
	var rlog *replicaLog
	if rlog, err = l.replicaLog(replica); err != nil {
		return nil, err
	}

	// Create the instance with default values and append it to leader's log
	inst = &pb.Instance{
		Replica: replica,
		Slot:    rlog.nextSlot(),
		Seq:     l.sequence,
		Deps:    make(map[uint32]uint64),
		Status:  pb.Status_INITIAL,
		Acks:    1,
		Ops:     ops,
	}

	// Sanity check that the append was succesful
	if err = rlog.insert(inst); err != nil {
		return nil, err
	}

	// Update dependencies and the conflict cache for the instance.
	l.updateDependencies(inst)
	l.updateConflicts(inst)

	return inst, nil
}

// Insert an instance into a log into the replica/slot specified by the instnace.
// Will return an error if there is already an instance in that slot or if inserting
// does not act like appending (e.g. the inst slot is not the next slot).
func (l *Logs) Insert(inst *pb.Instance) (err error) {
	var rlog *replicaLog
	if rlog, err = l.replicaLog(inst.Replica); err != nil {
		return err
	}

	return rlog.insert(inst)
}

// Helper function to insert an instance directly into a replica log.
func (l *replicaLog) insert(inst *pb.Instance) (err error) {
	next := l.nextSlot()
	if inst.Slot > next {
		return fmt.Errorf("cannot insert into slot %d when next slot is %d", inst.Slot, next)
	}

	if inst.Slot < next {
		if l.instances[inst.Slot] != nil {
			return fmt.Errorf("there is already an instance in slot %d", inst.Slot)
		}
		return fmt.Errorf("cannot append to slot %d before next slot %d", inst.Slot, next)
	}

	l.instances = append(l.instances, inst)
	return nil
}

// Get an instance in the specified replica's log at the specified index.
func (l *Logs) Get(replica uint32, slot uint64) (inst *pb.Instance, err error) {
	var rlog *replicaLog
	if rlog, err = l.replicaLog(replica); err != nil {
		return nil, err
	}

	if slot >= rlog.nextSlot() {
		return nil, fmt.Errorf("no instance found for replica PID %d in slot %d", replica, slot)
	}

	return rlog.instances[slot], nil
}

//===========================================================================
// Index Management (Slots, Dependencies, etc.)
//===========================================================================

// LastApplied returns the slot of the last applied instance for the specified replica.
func (l *Logs) LastApplied(replica uint32) (idx uint64, err error) {
	var rlog *replicaLog
	if rlog, err = l.replicaLog(replica); err != nil {
		return 0, err
	}

	next := rlog.nextSlot()
	if next == 0 {
		return 0, fmt.Errorf("no instances have been applied for replica PID %d", replica)
	}
	return next - 1, nil
}

// returns the next slot in the replica log.
func (l *replicaLog) nextSlot() uint64 {
	return uint64(len(l.instances))
}

// use the conflicts map to locate the latest dependency by slot across each replica's
// log and update the internal dependencies of the instance. Returns true if the
// dependencies on the instance have changed.
//
// TODO: should this simply happen on insert/append to the log?
func (l *Logs) updateDependencies(inst *pb.Instance) (changed bool) {
	// Ensure we have the latest dependency for all operations in the instance.
	for _, op := range inst.Ops {

		// Go through all replica logs to create the dependency map
		for pid, rlog := range l.logs {

			fmt.Println(rlog.conflicts)

			// If the replica has a conflict with this key add the conflict slot to the deps
			if slot, present := rlog.conflicts[op.Key]; present {

				fmt.Println(inst.Deps)

				// Check to see if the dependency has not changed
				if curdep, hasdep := inst.Deps[pid]; hasdep && slot <= curdep {
					// In this case the dependency is already stored or larger than the
					// cached conflict so no change in the deps needs to occur.
					continue
				}

				// Store the new depedency with the instance and mark changed
				// info("DEPENDENCY CHANGED! key %s pid %d slot %d", op.Key, pid, slot)
				changed = true
				inst.Deps[pid] = slot

				// Check to ensure that the instances sequence is bigger than all
				// conflicting instance sequences to ensure correct execution order.
				if rlog.instances[slot].Seq >= inst.Seq {
					inst.Seq = 1 + rlog.instances[slot].Seq
				}
			}
		}
	}

	// Ensure that our global sequence is monotonically increasing
	if inst.Seq > l.sequence {
		l.sequence = inst.Seq
		changed = true
	}

	return changed
}

// Update the local conflict cache for the leader's replica log by associating all the
// keys of the operations with the slot of the instance.
//
// TODO: should this simply happen on insert/append to the log?
func (l *Logs) updateConflicts(inst *pb.Instance) {
	// Get the replica log for the leader of the instance
	rlog := l.logs[inst.Replica]

	// Update the set of keys across all operations in the instance
	for _, op := range inst.Ops {
		// If the key is in the conflicts map, but there is already a conflict slot
		// greater than the slot of the instance, do not modify the conflicts log.
		if slot, present := rlog.conflicts[op.Key]; present && slot >= inst.Slot {
			continue
		}

		// Update the conflicts map if there was no conflict before or if a new conflict was added
		rlog.conflicts[op.Key] = inst.Slot
	}
}

//===========================================================================
// Helpers
//===========================================================================

func (l *Logs) replicaLog(replica uint32) (*replicaLog, error) {
	rlog, ok := l.logs[replica]
	if !ok {
		return nil, fmt.Errorf("no log for replica with PID %d", replica)
	}
	return rlog, nil
}
