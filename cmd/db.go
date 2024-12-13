package cmd

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/cryption"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"time"
)

func dbUp() {
	fmt.Println("Running database migrations...")
	var err error

	conn, closeConn := db.GetConnection()
	defer closeConn()

	// table
	err = conn.Migrator().AutoMigrate(&model.Photoinc{})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().AutoMigrate(&model.Photo{})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().AutoMigrate(&model.Company{})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().AutoMigrate(&model.Usercompany{})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().AutoMigrate(&model.Customer{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Order{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Paper{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Design{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Print{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Finishing{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Other{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Phase{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Orderphase{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Payment{})
	if err != nil {
		panic(err)
	}

	// view
	vPhoto := conn.Model(&model.Photo{}).Unscoped().
		Select("photos.*, u1.fullname as create_name").
		Joins("left join users u1 on u1.id = photos.create_by")

	err = conn.Migrator().CreateView(model.VIEW_PHOTO, gorm.ViewOption{
		Replace: true,
		Query:   vPhoto,
	})
	if err != nil {
		panic(err)
	}

	vUser := conn.Model(&model.User{}).Unscoped().
		Select("users.*, photos.photo_path as photo_url, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join photos photos on photos.id = users.photo_id").
		Joins("left join users u1 on u1.id = users.create_by").
		Joins("left join users u2 on u2.id = users.update_by")

	err = conn.Migrator().CreateView(model.VIEW_USER, gorm.ViewOption{
		Replace: true,
		Query:   vUser,
	})
	if err != nil {
		panic(err)
	}

	vCompany := conn.Model(&model.Company{}).Unscoped().
		Select("companies.*, photos.photo_path as photo_url, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join photos photos on photos.id = companies.photo_id").
		Joins("left join users u1 on u1.id = companies.create_by").
		Joins("left join users u2 on u2.id = companies.update_by")

	err = conn.Migrator().CreateView(model.VIEW_COMPANY, gorm.ViewOption{
		Replace: true,
		Query:   vCompany,
	})
	if err != nil {
		panic(err)
	}

	vUsercompany := conn.Model(&model.Usercompany{}).Unscoped().
		Select("usercompanies.*, companies.name as company_name, users.fullname as user_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = usercompanies.company_id").
		Joins("left join users users on users.id = usercompanies.user_id").
		Joins("left join users u1 on u1.id = usercompanies.create_by").
		Joins("left join users u2 on u2.id = usercompanies.update_by")

	err = conn.Migrator().CreateView(model.VIEW_USERCOMPANY, gorm.ViewOption{
		Replace: true,
		Query:   vUsercompany,
	})
	if err != nil {
		panic(err)
	}

	vCustomer := conn.Model(&model.Customer{}).Unscoped().
		Select("customers.*, companies.name as company_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = customers.company_id").
		Joins("left join users u1 on u1.id = customers.create_by").
		Joins("left join users u2 on u2.id = customers.update_by")

	err = conn.Migrator().CreateView(model.VIEW_CUSTOMER, gorm.ViewOption{
		Replace: true,
		Query:   vCustomer,
	})
	if err != nil {
		panic(err)
	}

	vOrder := conn.Model(&model.Order{}).Unscoped().
		Select("orders.*, companies.name as company_name, customers.name as customer_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = orders.company_id").
		Joins("left join customers customers on customers.id = orders.customer_id").
		Joins("left join users u1 on u1.id = orders.create_by").
		Joins("left join users u2 on u2.id = orders.update_by")

	err = conn.Migrator().CreateView(model.VIEW_ORDER, gorm.ViewOption{
		Replace: true,
		Query:   vOrder,
	})
	if err != nil {
		panic(err)
	}

	vPaper := conn.Model(&model.Paper{}).Unscoped().
		Select("papers.*, companies.name as company_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = papers.company_id").
		Joins("left join users u1 on u1.id = papers.create_by").
		Joins("left join users u2 on u2.id = papers.update_by")

	err = conn.Migrator().CreateView(model.VIEW_PAPER, gorm.ViewOption{
		Replace: true,
		Query:   vPaper,
	})
	if err != nil {
		panic(err)
	}

	vDesign := conn.Model(&model.Design{}).Unscoped().
		Select("designs.*, companies.name as company_name, designs.name as order_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = designs.company_id").
		Joins("left join orders orders on orders.id = designs.order_id").
		Joins("left join users u1 on u1.id = designs.create_by").
		Joins("left join users u2 on u2.id = designs.update_by")

	err = conn.Migrator().CreateView(model.VIEW_DESIGN, gorm.ViewOption{
		Replace: true,
		Query:   vDesign,
	})
	if err != nil {
		panic(err)
	}

	vPrint := conn.Model(&model.Print{}).Unscoped().
		Select("prints.*, companies.name as company_name, prints.name as print_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = prints.company_id").
		Joins("left join orders orders on orders.id = prints.order_id").
		Joins("left join users u1 on u1.id = prints.create_by").
		Joins("left join users u2 on u2.id = prints.update_by")

	err = conn.Migrator().CreateView(model.VIEW_PRINT, gorm.ViewOption{
		Replace: true,
		Query:   vPrint,
	})
	if err != nil {
		panic(err)
	}

	vFinishing := conn.Model(&model.Finishing{}).Unscoped().
		Select("finishings.*, companies.name as company_name, finishings.name as finishing_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = finishings.company_id").
		Joins("left join orders orders on orders.id = finishings.order_id").
		Joins("left join users u1 on u1.id = finishings.create_by").
		Joins("left join users u2 on u2.id = finishings.update_by")

	err = conn.Migrator().CreateView(model.VIEW_FINISHING, gorm.ViewOption{
		Replace: true,
		Query:   vFinishing,
	})
	if err != nil {
		panic(err)
	}

	vOther := conn.Model(&model.Other{}).Unscoped().
		Select("others.*, companies.name as company_name, orders.name as order_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = others.company_id").
		Joins("left join orders orders on orders.id = others.order_id").
		Joins("left join users u1 on u1.id = others.create_by").
		Joins("left join users u2 on u2.id = others.update_by")

	err = conn.Migrator().CreateView(model.VIEW_OTHER, gorm.ViewOption{
		Replace: true,
		Query:   vOther,
	})
	if err != nil {
		panic(err)
	}

	vPhase := conn.Model(&model.Phase{}).Unscoped().
		Select("phases.*, companies.name as company_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = phases.company_id").
		Joins("left join users u1 on u1.id = phases.create_by").
		Joins("left join users u2 on u2.id = phases.update_by")

	err = conn.Migrator().CreateView(model.VIEW_PHASE, gorm.ViewOption{
		Replace: true,
		Query:   vPhase,
	})
	if err != nil {
		panic(err)
	}

	vOrderphase := conn.Model(&model.Orderphase{}).Unscoped().
		Select("orderphases.*, companies.name as company_name, orderphases.name as orderphase_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = orderphases.company_id").
		Joins("left join orders orders on orders.id = orderphases.order_id").
		Joins("left join users u1 on u1.id = orderphases.create_by").
		Joins("left join users u2 on u2.id = orderphases.update_by")

	err = conn.Migrator().CreateView(model.VIEW_ORDERPHASE, gorm.ViewOption{
		Replace: true,
		Query:   vOrderphase,
	})
	if err != nil {
		panic(err)
	}

	vPayment := conn.Model(&model.Payment{}).Unscoped().
		Select("payments.*, companies.name as company_name, payments.name as payment_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = payments.company_id").
		Joins("left join orders orders on orders.id = payments.order_id").
		Joins("left join users u1 on u1.id = payments.create_by").
		Joins("left join users u2 on u2.id = payments.update_by")

	err = conn.Migrator().CreateView(model.VIEW_PAYMENT, gorm.ViewOption{
		Replace: true,
		Query:   vPayment,
	})
	if err != nil {
		panic(err)
	}

}

func dbDown() {
	fmt.Println("Reverting database migrations...")
	var err error

	conn, closeConn := db.GetConnection()
	defer closeConn()

	err = conn.Exec("DROP SCHEMA public CASCADE").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("CREATE SCHEMA public").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("GRANT ALL ON SCHEMA public TO postgres").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("GRANT ALL ON SCHEMA public TO public").Error
	if err != nil {
		panic(err)
	}
}

func dbSeed() {
	fmt.Println("Seeding the database with initial data start")

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	userID := utils.GetUniqueID()
	demoUserID := utils.GetUniqueID()
	demoCompanyID := utils.GetUniqueID()

	now := time.Now()

	password, err := cryption.EncryptAES64("123456")
	if err != nil {
		panic(err)
	}

	users := []model.User{
		{ID: userID, Role: constant.RoleAdmin, Email: "jihanlugas2@gmail.com", Username: "jihanlugas", NoHp: utils.FormatPhoneTo62("6287770333043"), Fullname: "Jihan Lugas", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now, CreateBy: userID, UpdateBy: userID},
		{ID: demoUserID, Role: constant.RoleUseradmin, Email: "admindemo@gmail.com", Username: "admindemo", NoHp: utils.FormatPhoneTo62("6287770331234"), Fullname: "Admin Demo", Passwd: password, PassVersion: 1, IsActive: true, AccountVerifiedDt: &now, CreateBy: userID, UpdateBy: userID},
	}
	tx.Create(&users)

	companies := []model.Company{
		{ID: demoCompanyID, Name: "Demo Company", Description: "Demo Company Generated", CreateBy: userID, UpdateBy: userID},
	}
	tx.Create(&companies)

	usercompanies := []model.Usercompany{
		{
			UserID:           demoUserID,
			CompanyID:        demoCompanyID,
			IsDefaultCompany: true,
			IsCreator:        true,
			CreateBy:         userID,
			UpdateBy:         userID,
		},
	}

	tx.Create(&usercompanies)

	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}

	fmt.Println("Seeding the database with initial data end")
}

func dbReset() {
	dbDown()
	dbUp()
	dbSeed()
}
