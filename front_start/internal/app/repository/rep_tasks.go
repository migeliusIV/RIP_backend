package repository

/*
// GetOrCreateDraftFrax находит заявку-черновик для пользователя или создает новую.
// Временно используем userID = 1, как "захардкоженный" ID пользователя-модератора.
func (r *Repository) GetOrCreateDraftFrax(userID uint) (*ds.FraxSearching, error) {
	var frax ds.FraxSearching

	// Ищем черновик у пользователя
	err := r.db.Where("creator_id = ? AND status = ?", userID, ds.StatusDraft).First(&frax).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Если не нашли, создаем новый черновик
		newFrax := ds.FraxSearching{
			CreatorID: userID,
			Status:    ds.StatusDraft,
		}
		if err := r.db.Create(&newFrax).Error; err != nil {
			return nil, err
		}
		return &newFrax, nil
	}

	return &frax, err
}

// AddFactorToFrax добавляет фактор в заявку (создает запись в таблице m-m).
func (r *Repository) AddFactorToFrax(fraxID, factorID uint) error {
	// Проверяем, нет ли уже такого фактора в заявке, чтобы избежать дублей
	var count int64
	r.db.Model(&ds.FactorToFrax{}).Where("frax_id = ? AND factor_id = ?", fraxID, factorID).Count(&count)
	if count > 0 {
		return errors.New("factor already in frax")
	}

	link := ds.FactorToFrax{
		FraxID:   fraxID,
		FactorID: factorID,
	}
	return r.db.Create(&link).Error
}

// GetFraxWithFactors получает заявку со всеми связанными факторами.
// Используем Preload для эффективной загрузки связанных данных.
func (r *Repository) GetFraxWithFactors(fraxID uint) (*ds.FraxSearching, error) {
	var frax ds.FraxSearching

	err := r.db.Preload("FactorsLink.Factor").First(&frax, fraxID).Error
	if err != nil {
		return nil, err
	}

	// Проверяем, что заявка не удалена
	if frax.Status == ds.StatusDeleted {
		return nil, errors.New("frax page not found or has been deleted")
	}

	return &frax, nil
}

// LogicallyDeleteFrax выполняет логическое удаление заявки через чистый SQL UPDATE.
func (r *Repository) LogicallyDeleteFrax(fraxID uint) error {
	// Используем Exec для выполнения "сырого" SQL-запроса
	result := r.db.Exec("UPDATE frax_searchings SET status = ? WHERE id = ?", ds.StatusDeleted, fraxID)
	return result.Error
}
*/
