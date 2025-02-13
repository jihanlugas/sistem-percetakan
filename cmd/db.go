package cmd

import (
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/cryption"
	"github.com/jihanlugas/sistem-percetakan/db"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

func dbUp() {
	fmt.Println("Running database migrations...")
	dbUpTable()
	dbUpView()
}

func dbUpTable() {
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
}

func dbUpView() {
	var err error

	conn, closeConn := db.GetConnection()
	defer closeConn()

	// view
	err = conn.Migrator().DropView(model.VIEW_PHOTO)
	if err != nil {
		panic(err)
	}
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

	err = conn.Migrator().DropView(model.VIEW_USER)
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

	err = conn.Migrator().DropView(model.VIEW_COMPANY)
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

	err = conn.Migrator().DropView(model.VIEW_USERCOMPANY)
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

	err = conn.Migrator().DropView(model.VIEW_CUSTOMER)
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

	err = conn.Migrator().DropView(model.VIEW_ORDER)
	if err != nil {
		panic(err)
	}
	vOrder := conn.Model(&model.Order{}).Unscoped().
		Select("orders.*" +
			", orderphases.id as orderphase_id" +
			", orderphases.phase_id as phase_id" +
			", orderphases.name as orderphase_name" +
			", coalesce(designs.total_design, 0) as total_design" +
			", coalesce(prints.total_print, 0) as total_print" +
			", coalesce(finishings.total_finishing, 0) as total_finishing" +
			", coalesce(others.total_other, 0) as total_other" +
			", coalesce(payments.total_payment, 0) as total_payment" +
			", coalesce(designs.total_design, 0) + coalesce(prints.total_print, 0) + coalesce(finishings.total_finishing, 0) + coalesce(others.total_other, 0) as total_order" +
			", coalesce(designs.total_design, 0) + coalesce(prints.total_print, 0) + coalesce(finishings.total_finishing, 0) + coalesce(others.total_other, 0) - coalesce(payments.total_payment, 0) as outstanding" +
			", companies.name as company_name, customers.name as customer_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join orderphases orderphases on orderphases.order_id = orders.id " +
			"AND orderphases.create_dt = (select max(orderphases.create_dt) from orderphases where orderphases.order_id = orders.id) " +
			"AND orderphases.delete_dt is null").
		Joins("left join ( " +
			"select d.order_id, COALESCE(sum(d.total), 0) as total_design " +
			"from designs d " +
			"where d.delete_dt is null " +
			"group by d.order_id " +
			") as designs on designs.order_id = orders.id").
		Joins("left join ( " +
			"select p.order_id, COALESCE(sum(p.total), 0) as total_print " +
			"from prints p " +
			"where p.delete_dt is null " +
			"group by p.order_id " +
			") as prints on prints.order_id = orders.id").
		Joins("left join ( " +
			"select f.order_id, COALESCE(sum(f.total), 0) as total_finishing " +
			"from finishings f " +
			"where f.delete_dt is null " +
			"group by f.order_id " +
			") as finishings on finishings.order_id = orders.id").
		Joins("left join ( " +
			"select o.order_id, COALESCE(sum(o.total), 0) as total_other " +
			"from others o " +
			"where o.delete_dt is null " +
			"group by o.order_id " +
			") as others on others.order_id = orders.id").
		Joins("left join ( " +
			"select p.order_id, COALESCE(sum(p.amount), 0) as total_payment " +
			"from payments p " +
			"where p.delete_dt is null " +
			"group by p.order_id " +
			") as payments on payments.order_id = orders.id").
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

	err = conn.Migrator().DropView(model.VIEW_PAPER)
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

	err = conn.Migrator().DropView(model.VIEW_DESIGN)
	if err != nil {
		panic(err)
	}
	vDesign := conn.Model(&model.Design{}).Unscoped().
		Select("designs.*, companies.name as company_name, orders.name as order_name, u1.fullname as create_name, u2.fullname as update_name").
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

	err = conn.Migrator().DropView(model.VIEW_PRINT)
	if err != nil {
		panic(err)
	}
	vPrint := conn.Model(&model.Print{}).Unscoped().
		Select("prints.*, companies.name as company_name, orders.name as order_name, papers.name as paper_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = prints.company_id").
		Joins("left join orders orders on orders.id = prints.order_id").
		Joins("left join papers papers on papers.id = prints.paper_id").
		Joins("left join users u1 on u1.id = prints.create_by").
		Joins("left join users u2 on u2.id = prints.update_by")
	err = conn.Migrator().CreateView(model.VIEW_PRINT, gorm.ViewOption{
		Replace: true,
		Query:   vPrint,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_FINISHING)
	if err != nil {
		panic(err)
	}
	vFinishing := conn.Model(&model.Finishing{}).Unscoped().
		Select("finishings.*, companies.name as company_name, orders.name as order_name, u1.fullname as create_name, u2.fullname as update_name").
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

	err = conn.Migrator().DropView(model.VIEW_OTHER)
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

	err = conn.Migrator().DropView(model.VIEW_PHASE)
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

	err = conn.Migrator().DropView(model.VIEW_ORDERPHASE)
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

	err = conn.Migrator().DropView(model.VIEW_PAYMENT)
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
		{
			ID:                userID,
			Role:              constant.RoleAdmin,
			Email:             "jihanlugas2@gmail.com",
			Username:          "jihanlugas",
			PhoneNumber:       utils.FormatPhoneTo62("6287770333043"),
			Fullname:          "Jihan Lugas",
			Address:           "Jl. Gunung Sahari No. 10, Jakarta Pusat",
			Passwd:            password,
			PassVersion:       1,
			IsActive:          true,
			AccountVerifiedDt: &now,
			CreateBy:          userID,
			UpdateBy:          userID},
		{
			ID:                demoUserID,
			Role:              constant.RoleUseradmin,
			Email:             "admindemo@gmail.com",
			Username:          "admindemo",
			PhoneNumber:       utils.FormatPhoneTo62("6287770331234"),
			Fullname:          "Admin Demo",
			Address:           "Jl. Raya Jatinegara No. 10, Jakarta Timur",
			Passwd:            password,
			PassVersion:       1,
			IsActive:          true,
			AccountVerifiedDt: &now,
			CreateBy:          userID,
			UpdateBy:          userID,
		},
	}
	tx.Create(&users)

	companies := []model.Company{
		{ID: demoCompanyID,
			Name:        "Demo Company",
			Description: "Demo Company Generated",
			Email:       "companydemo@gmail",
			PhoneNumber: utils.FormatPhoneTo62("6287770331234"),
			Address:     "Jl. M.H. Thamrin No. 10, Jakarta Pusat",
			CreateBy:    userID,
			UpdateBy:    userID,
		},
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

	phases := []model.Phase{
		{ID: utils.GetUniqueID(), Order: 1, Name: "Open", Description: "Order Masuk", CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 2, Name: "Design", Description: "Order sedang di desain", CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 3, Name: "Print", Description: "Order sedang di print", CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 4, Name: "Finishing", Description: "Order sedang di finishing", CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 5, Name: "Done", Description: "Order selesai proses", CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 6, Name: "Close", Description: "Order sudah di berikan customer", CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
	}
	tx.Create(&phases)

	papers := []model.Paper{
		{ID: utils.GetUniqueID(), Name: "A3 HVS 80", Description: "", DefaultPrice: 1800, DefaultPriceDuplex: 2600, CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 HVS 100", Description: "", DefaultPrice: 1900, DefaultPriceDuplex: 2700, CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Paper 120", Description: "", DefaultPrice: 2000, DefaultPriceDuplex: 2800, CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Paper 150", Description: "", DefaultPrice: 2100, DefaultPriceDuplex: 2900, CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Carton 210", Description: "", DefaultPrice: 2200, DefaultPriceDuplex: 3000, CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Carton 230", Description: "", DefaultPrice: 2300, DefaultPriceDuplex: 3100, CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Carton 260", Description: "", DefaultPrice: 2400, DefaultPriceDuplex: 3200, CompanyID: demoCompanyID, CreateBy: userID, UpdateBy: userID},
	}
	tx.Create(&papers)

	customers := []model.Customer{
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Ikuta Rira", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Udin", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Budi", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Roronoa Zoro", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Kamado Tanjiro", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Kim Jong Un", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Uru", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Al Ghazali", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Aliando", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Alguero", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Cahyono", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Cristiano Ronaldo", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Leonel Messi", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Gojo Satoru", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Park Chan-wook", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: demoCompanyID, Name: "Lee Myung-bak", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
	}
	mapCustomers := []model.Customer{}
	for i, customer := range customers {
		mapCustomer := customer
		mapCustomer.Email = fmt.Sprintf("%s@gmail.com", strings.ToLower(strings.ReplaceAll(customer.Name, " ", "")))
		mapCustomer.Address = fmt.Sprintf("%s %d", mapCustomer.Address, i)
		mapCustomers = append(mapCustomers, mapCustomer)
	}
	tx.Create(&mapCustomers)

	orders := []model.Order{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			CustomerID:  customers[0].ID,
			Name:        "Majalah Manchester City",
			Description: "- Majalah Liga Inggris\n- Majalah Liga Champions",
			IsDone:      false,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			CustomerID:  customers[1].ID,
			Name:        "Batik",
			Description: "1. 200 lembar\n2. 300 lembar\n3. 50 lembar",
			IsDone:      true,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
	}
	tx.Create(&orders)

	designs := []model.Design{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			Name:        "Cover Liga Inggris",
			Description: "Photo photo liga inggris",
			Qty:         1,
			Price:       20000,
			Total:       1 * 20000,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			Name:        "Cover iga champions",
			Description: "Photo liga champions",
			Qty:         1,
			Price:       30000,
			Total:       1 * 30000,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
	}
	tx.Create(&designs)

	prints := []model.Print{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			PaperID:     papers[5].ID,
			Name:        "Cover Liga Inggris",
			Description: "Cover 400 lembar",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         400,
			Price:       2300,
			Total:       1 * 400 * 2300,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			PaperID:     papers[6].ID,
			Name:        "Cover Liga Champions",
			Description: "Cover 300 lembar",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         300,
			Price:       2400,
			Total:       1 * 300 * 2400,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			PaperID:     papers[2].ID,
			Name:        "Isi Liga Inggris",
			Description: "Isi timbal balik 120 lembar",
			IsDuplex:    true,
			PageCount:   120,
			Qty:         400,
			Price:       2800,
			Total:       120 * 400 * 2800,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			PaperID:     papers[3].ID,
			Name:        "Isi Liga Champions",
			Description: "Isi timbal balik 140 lembar",
			IsDuplex:    true,
			PageCount:   140,
			Qty:         300,
			Price:       2900,
			Total:       140 * 300 * 2900,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[1].ID,
			PaperID:     papers[6].ID,
			Name:        "Model 1",
			Description: "300 Lembar",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         300,
			Price:       2400,
			Total:       1 * 300 * 2400,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[1].ID,
			PaperID:     papers[6].ID,
			Name:        "Model 2",
			Description: "400 Lembar",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         400,
			Price:       2400,
			Total:       1 * 400 * 2400,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[1].ID,
			PaperID:     papers[6].ID,
			Name:        "Model 3",
			Description: "50 Lembar",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         50,
			Price:       2400,
			Total:       1 * 50 * 2400,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
	}
	tx.Create(&prints)

	finishings := []model.Finishing{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			Name:        "Potong",
			Description: "Borongan potong 200K",
			Qty:         1,
			Price:       200000,
			Total:       1 * 200000,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			Name:        "Laminating Cover",
			Description: "Laminating Cover 300 dan 400",
			Qty:         700,
			Price:       2400,
			Total:       700 * 2400,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			Name:        "Banding",
			Description: "Banding 300 dan 400 Buku",
			Qty:         700,
			Price:       17000,
			Total:       700 * 17000,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
	}
	tx.Create(&finishings)

	others := []model.Other{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			Name:        "Sampul liga inggris",
			Description: "Borongan ke belakang",
			Qty:         400,
			Price:       18000,
			Total:       400 * 18000,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[0].ID,
			Name:        "Sampul liga champion",
			Description: "Borongan ke belakang",
			Qty:         300,
			Price:       20000,
			Total:       300 * 20000,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
	}
	tx.Create(&others)

	payments := []model.Payment{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[1].ID,
			Name:        "DP",
			Description: "DP 10 jt",
			Amount:      1000000,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   demoCompanyID,
			OrderID:     orders[1].ID,
			Name:        "Lunas",
			Description: "Lunas",
			Amount:      800000,
			CreateBy:    userID,
			UpdateBy:    userID,
		},
	}
	tx.Create(&payments)

	orderphases := []model.Orderphase{
		{
			ID:        utils.GetUniqueID(),
			CompanyID: demoCompanyID,
			OrderID:   orders[0].ID,
			PhaseID:   phases[0].ID,
			Name:      phases[0].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, -3),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, -3),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: demoCompanyID,
			OrderID:   orders[0].ID,
			PhaseID:   phases[1].ID,
			Name:      phases[1].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, -2),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, -2),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: demoCompanyID,
			OrderID:   orders[1].ID,
			PhaseID:   phases[0].ID,
			Name:      phases[0].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, -3),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, -3),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: demoCompanyID,
			OrderID:   orders[1].ID,
			PhaseID:   phases[1].ID,
			Name:      phases[1].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, -2),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, -2),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: demoCompanyID,
			OrderID:   orders[1].ID,
			PhaseID:   phases[2].ID,
			Name:      phases[2].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, -1),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, -1),
		},
	}
	tx.Create(&orderphases)

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
