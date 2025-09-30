package repository

import (
    "errors"
    "front_start/internal/app/ds"
)

func (r *Repository) RegisterUser(login, password string) error {
    u := ds.Users{Login: login, Password: password}
    return r.db.Create(&u).Error
}

func (r *Repository) GetUserByID(id uint) (*ds.Users, error) {
    var u ds.Users
    if err := r.db.First(&u, "id_user = ?", id).Error; err != nil {
        return nil, err
    }
    u.Password = ""
    return &u, nil
}

func (r *Repository) UpdateUser(id uint, password *string) (*ds.Users, error) {
    var u ds.Users
    if err := r.db.First(&u, "id_user = ?", id).Error; err != nil {
        return nil, err
    }
    if password != nil && *password != "" {
        u.Password = *password
    }
    if err := r.db.Save(&u).Error; err != nil {
        return nil, err
    }
    u.Password = ""
    return &u, nil
}

func (r *Repository) CheckUser(login, password string) error {
    var u ds.Users
    if err := r.db.Where("login = ?", login).First(&u).Error; err != nil {
        return err
    }
    if u.Password != password {
        return errors.New("invalid credentials")
    }
    return nil
}


