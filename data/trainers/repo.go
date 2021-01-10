package trainers

import "github.com/AlexanderYukhanov/90minsGOexp/data/models"

// Repo - trainers repo.
type Repo interface {
	FindTrainer(trainer models.TrainerID) (*models.Trainer, error)
}
