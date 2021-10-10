package sequencer

import "github.com/daginwu/api/slog/pb"

type Sequencer struct {
	config Config
}

type Config struct {
	BatchSize int
}

type LocalLog struct {
}

type GlobalLog struct {
}

func InitSequencer(config Config) *Sequencer {
	return &Sequencer{
		config: config,
	}
}

func AnalyzeTransaction(txn *pb.Transaction) *pb.Transaction {

	return txn
}
