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
	err = conn.Migrator().AutoMigrate(&model.Print{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Finishing{})
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
	err = conn.Migrator().AutoMigrate(&model.Transaction{})
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
		Select("users.*, usercompanies.id as usercompany_id, usercompanies.company_id as company_id, '' as photo_url, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join usercompanies usercompanies on usercompanies.user_id = users.id").
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
		Select("companies.*, '' as photo_url, u1.fullname as create_name, u2.fullname as update_name").
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
			", coalesce(prints.total_print, 0) as total_print" +
			", coalesce(finishings.total_finishing, 0) as total_finishing" +
			", coalesce(transactions.total_transaction, 0) as total_transaction" +
			", coalesce(prints.total_print, 0) + coalesce(finishings.total_finishing, 0) as total_order" +
			", coalesce(prints.total_print, 0) + coalesce(finishings.total_finishing, 0) - coalesce(transactions.total_transaction, 0) as outstanding" +
			", companies.name as company_name, customers.name as customer_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join orderphases orderphases on orderphases.order_id = orders.id " +
			"AND orderphases.create_dt = (select max(orderphases.create_dt) from orderphases where orderphases.order_id = orders.id) " +
			"AND orderphases.delete_dt is null").
		Joins("left join ( " +
			"select p.order_id, COALESCE(sum(p.total), 0) as total_print " +
			"from prints p " +
			"where p.delete_dt is null " +
			"group by p.order_id " +
			") as prints on prints.order_id = orders.id").
		Joins("left join ( " +
			"select o.order_id, COALESCE(sum(o.total), 0) as total_finishing " +
			"from finishings o " +
			"where o.delete_dt is null " +
			"group by o.order_id " +
			") as finishings on finishings.order_id = orders.id").
		Joins("left join ( " +
			"select p.order_id, COALESCE(sum(p.amount), 0) as total_transaction " +
			"from transactions p " +
			"where p.delete_dt is null " +
			"group by p.order_id " +
			") as transactions on transactions.order_id = orders.id").
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
		Select("orderphases.*, companies.name as company_name, orders.name as order_name, u1.fullname as create_name, u2.fullname as update_name").
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

	err = conn.Migrator().DropView(model.VIEW_TRANSACTION)
	if err != nil {
		panic(err)
	}
	vTransaction := conn.Model(&model.Transaction{}).Unscoped().
		Select("transactions.*, companies.name as company_name, orders.name as order_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join companies companies on companies.id = transactions.company_id").
		Joins("left join orders orders on orders.id = transactions.order_id").
		Joins("left join users u1 on u1.id = transactions.create_by").
		Joins("left join users u2 on u2.id = transactions.update_by")

	err = conn.Migrator().CreateView(model.VIEW_TRANSACTION, gorm.ViewOption{
		Replace: true,
		Query:   vTransaction,
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

	adminID := utils.GetUniqueID()
	userID := "f7416f17-884b-46d3-b7db-b90be60a71c5"
	companyID := "fcc18dfc-b0ef-42ef-8036-28503492a2a1"
	companyPhotoID := "05421e10-6e09-4ae0-9102-7947b2166d30"
	
	now := time.Now()

	password, err := cryption.EncryptAES64("123456")
	if err != nil {
		panic(err)
	}

	users := []model.User{
		{
			ID:                adminID,
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
			CreateBy:          adminID,
			UpdateBy:          adminID,
		},
		{
			ID:                userID,
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
			CreateBy:          adminID,
			UpdateBy:          adminID,
		},
	}
	tx.Create(&users)

	companies := []model.Company{
		{
			ID:          companyID,
			Name:        "Demo Company",
			Description: "Demo Company Generated",
			Email:       "companydemo@gmail",
			PhoneNumber: utils.FormatPhoneTo62("6287770331234"),
			Address:     "Jl. M.H. Thamrin No. 10, Jakarta Pusat",
			PhotoID:     companyPhotoID,
			CreateBy:    adminID,
			UpdateBy:    adminID,
		},
	}
	tx.Create(&companies)

	usercompanies := []model.Usercompany{
		{
			UserID:           userID,
			CompanyID:        companyID,
			IsDefaultCompany: true,
			IsCreator:        true,
			CreateBy:         adminID,
			UpdateBy:         adminID,
		},
	}
	tx.Create(&usercompanies)

	photos := []model.Photo{
		{
			ID:         companyPhotoID,
			ClientName: "Logo",
			ServerName: "Logo",
			RefTable:   "",
			Ext:        "png",
			PhotoPath:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAlgAAAJYCAYAAAC+ZpjcAAAACXBIWXMAAC4jAAAuIwF4pT92AAAJymlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4gPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iQWRvYmUgWE1QIENvcmUgNS42LWMxNDggNzkuMTY0MDM2LCAyMDE5LzA4LzEzLTAxOjA2OjU3ICAgICAgICAiPiA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtbG5zOmRjPSJodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyIgeG1sbnM6cGhvdG9zaG9wPSJodHRwOi8vbnMuYWRvYmUuY29tL3Bob3Rvc2hvcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RFdnQ9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZUV2ZW50IyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIDIzLjQgKFdpbmRvd3MpIiB4bXA6Q3JlYXRlRGF0ZT0iMjAyNC0wMi0xOVQxMzoxMjowNS0wODowMCIgeG1wOk1vZGlmeURhdGU9IjIwMjUtMDItMjJUMTM6MzM6MDkrMDc6MDAiIHhtcDpNZXRhZGF0YURhdGU9IjIwMjUtMDItMjJUMTM6MzM6MDkrMDc6MDAiIGRjOmZvcm1hdD0iaW1hZ2UvcG5nIiBwaG90b3Nob3A6Q29sb3JNb2RlPSIzIiB4bXBNTTpJbnN0YW5jZUlEPSJ4bXAuaWlkOmFlYzYxMDY4LWFjMmQtOWE0Yy05YmVlLWZjZjFmY2Q4ZjdkYiIgeG1wTU06RG9jdW1lbnRJRD0iYWRvYmU6ZG9jaWQ6cGhvdG9zaG9wOmY5ZjQ4N2Y2LTk2MmYtNjk0My05MmE0LWNhNjM1MDM4YTU0NSIgeG1wTU06T3JpZ2luYWxEb2N1bWVudElEPSJ4bXAuZGlkOjAxNmU4MjI3LTJkMDItNDE0Mi04NjY2LWQzMGNmMDhkZTBmYyI+IDx4bXBNTTpIaXN0b3J5PiA8cmRmOlNlcT4gPHJkZjpsaSBzdEV2dDphY3Rpb249ImNyZWF0ZWQiIHN0RXZ0Omluc3RhbmNlSUQ9InhtcC5paWQ6MDE2ZTgyMjctMmQwMi00MTQyLTg2NjYtZDMwY2YwOGRlMGZjIiBzdEV2dDp3aGVuPSIyMDI0LTAyLTE5VDEzOjEyOjA1LTA4OjAwIiBzdEV2dDpzb2Z0d2FyZUFnZW50PSJBZG9iZSBQaG90b3Nob3AgMjMuNCAoV2luZG93cykiLz4gPHJkZjpsaSBzdEV2dDphY3Rpb249ImNvbnZlcnRlZCIgc3RFdnQ6cGFyYW1ldGVycz0iZnJvbSBpbWFnZS9wbmcgdG8gYXBwbGljYXRpb24vdm5kLmFkb2JlLnBob3Rvc2hvcCIvPiA8cmRmOmxpIHN0RXZ0OmFjdGlvbj0ic2F2ZWQiIHN0RXZ0Omluc3RhbmNlSUQ9InhtcC5paWQ6YWVlYjYyNjItY2VjYS1mMjRiLTg4YzQtMTMyOTI2MGE2NjEwIiBzdEV2dDp3aGVuPSIyMDI0LTA4LTMwVDA4OjAzOjA1LTA3OjAwIiBzdEV2dDpzb2Z0d2FyZUFnZW50PSJBZG9iZSBQaG90b3Nob3AgMjMuNCAoV2luZG93cykiIHN0RXZ0OmNoYW5nZWQ9Ii8iLz4gPHJkZjpsaSBzdEV2dDphY3Rpb249InNhdmVkIiBzdEV2dDppbnN0YW5jZUlEPSJ4bXAuaWlkOjMwZWU0N2I1LTFiNWUtN2E0Ny05OWYzLWY1NDNiZjE1MjYxYyIgc3RFdnQ6d2hlbj0iMjAyNS0wMi0yMlQxMzozMzowOSswNzowMCIgc3RFdnQ6c29mdHdhcmVBZ2VudD0iQWRvYmUgUGhvdG9zaG9wIDIxLjAgKFdpbmRvd3MpIiBzdEV2dDpjaGFuZ2VkPSIvIi8+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJjb252ZXJ0ZWQiIHN0RXZ0OnBhcmFtZXRlcnM9ImZyb20gYXBwbGljYXRpb24vdm5kLmFkb2JlLnBob3Rvc2hvcCB0byBpbWFnZS9wbmciLz4gPHJkZjpsaSBzdEV2dDphY3Rpb249ImRlcml2ZWQiIHN0RXZ0OnBhcmFtZXRlcnM9ImNvbnZlcnRlZCBmcm9tIGFwcGxpY2F0aW9uL3ZuZC5hZG9iZS5waG90b3Nob3AgdG8gaW1hZ2UvcG5nIi8+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJzYXZlZCIgc3RFdnQ6aW5zdGFuY2VJRD0ieG1wLmlpZDphZWM2MTA2OC1hYzJkLTlhNGMtOWJlZS1mY2YxZmNkOGY3ZGIiIHN0RXZ0OndoZW49IjIwMjUtMDItMjJUMTM6MzM6MDkrMDc6MDAiIHN0RXZ0OnNvZnR3YXJlQWdlbnQ9IkFkb2JlIFBob3Rvc2hvcCAyMS4wIChXaW5kb3dzKSIgc3RFdnQ6Y2hhbmdlZD0iLyIvPiA8L3JkZjpTZXE+IDwveG1wTU06SGlzdG9yeT4gPHhtcE1NOkRlcml2ZWRGcm9tIHN0UmVmOmluc3RhbmNlSUQ9InhtcC5paWQ6MzBlZTQ3YjUtMWI1ZS03YTQ3LTk5ZjMtZjU0M2JmMTUyNjFjIiBzdFJlZjpkb2N1bWVudElEPSJhZG9iZTpkb2NpZDpwaG90b3Nob3A6NjZlNzY3ZmMtMTk0Yy1hNDRhLTkyZDItYTFjMzBlMTQxNDJjIiBzdFJlZjpvcmlnaW5hbERvY3VtZW50SUQ9InhtcC5kaWQ6MDE2ZTgyMjctMmQwMi00MTQyLTg2NjYtZDMwY2YwOGRlMGZjIi8+IDwvcmRmOkRlc2NyaXB0aW9uPiA8L3JkZjpSREY+IDwveDp4bXBtZXRhPiA8P3hwYWNrZXQgZW5kPSJyIj8+dtJWmwABCrxJREFUeJzs/eeXJPl13w1+7i8iM8u0N9Pd0zM93hsMMBhYAiAIkRBoBRrRipQh9cid3X2x5+w5u3/A7jn74jlHj6RHK5EypCiJFiJBggThgQGIgRtvu2d6elx7XyYzI+J398WNyMyqyiybVZXVfT9zcqo6TURkVmTEN675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",
			CreateBy:   userID,
		},
	}
	tx.Create(&photos)

	phases := []model.Phase{
		{ID: utils.GetUniqueID(), Order: 1, Name: "Open", Description: "Order Masuk", CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 2, Name: "Design", Description: "Order sedang di desain", CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 3, Name: "Print", Description: "Order sedang di print", CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 4, Name: "Finishing", Description: "Order sedang di finishing", CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 5, Name: "Done", Description: "Order selesai proses", CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Order: 6, Name: "Close", Description: "Order sudah di berikan customer", CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
	}
	tx.Create(&phases)

	papers := []model.Paper{
		{ID: utils.GetUniqueID(), Name: "A3 HVS 80", Description: "", DefaultPrice: 1800, DefaultPriceDuplex: 2600, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 HVS 100", Description: "", DefaultPrice: 1900, DefaultPriceDuplex: 2700, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Paper 120", Description: "", DefaultPrice: 2000, DefaultPriceDuplex: 2800, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Paper 150", Description: "", DefaultPrice: 2100, DefaultPriceDuplex: 2900, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Carton 210", Description: "", DefaultPrice: 2200, DefaultPriceDuplex: 3000, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Carton 230", Description: "", DefaultPrice: 2300, DefaultPriceDuplex: 3100, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Art Carton 260", Description: "", DefaultPrice: 2400, DefaultPriceDuplex: 3200, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Stiker Glosy", Description: "", DefaultPrice: 5000, DefaultPriceDuplex: 7000, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A3 Stiker Dof", Description: "", DefaultPrice: 6000, DefaultPriceDuplex: 8000, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 HVS 80", Description: "", DefaultPrice: 900, DefaultPriceDuplex: 1300, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 HVS 100", Description: "", DefaultPrice: 950, DefaultPriceDuplex: 1350, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 Art Paper 120", Description: "", DefaultPrice: 1000, DefaultPriceDuplex: 1400, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 Art Paper 150", Description: "", DefaultPrice: 1050, DefaultPriceDuplex: 1450, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 Art Carton 210", Description: "", DefaultPrice: 1100, DefaultPriceDuplex: 1500, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 Art Carton 230", Description: "", DefaultPrice: 1150, DefaultPriceDuplex: 1550, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 Art Carton 260", Description: "", DefaultPrice: 1200, DefaultPriceDuplex: 1600, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 Stiker Glosy", Description: "", DefaultPrice: 2500, DefaultPriceDuplex: 3500, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A4 Stiker Dof", Description: "", DefaultPrice: 3000, DefaultPriceDuplex: 4000, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 HVS 80", Description: "", DefaultPrice: 450, DefaultPriceDuplex: 650, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 HVS 100", Description: "", DefaultPrice: 475, DefaultPriceDuplex: 675, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 Art Paper 120", Description: "", DefaultPrice: 500, DefaultPriceDuplex: 700, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 Art Paper 150", Description: "", DefaultPrice: 525, DefaultPriceDuplex: 725, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 Art Carton 210", Description: "", DefaultPrice: 550, DefaultPriceDuplex: 750, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 Art Carton 230", Description: "", DefaultPrice: 575, DefaultPriceDuplex: 775, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 Art Carton 260", Description: "", DefaultPrice: 600, DefaultPriceDuplex: 800, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 Stiker Glosy", Description: "", DefaultPrice: 1250, DefaultPriceDuplex: 1750, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "A5 Stiker Dof", Description: "", DefaultPrice: 1500, DefaultPriceDuplex: 2000, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), Name: "Blangko", Description: "", DefaultPrice: 5000, DefaultPriceDuplex: 7000, CompanyID: companyID, CreateBy: userID, UpdateBy: userID},
	}
	tx.Create(&papers)

	customers := []model.Customer{
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Ikuta Rira", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Udin", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Budi", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Roronoa Zoro", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Kamado Tanjiro", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Kim Jong Un", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Uru", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Al Ghazali", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Aliando", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Alguero", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Cahyono", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Cristiano Ronaldo", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Leonel Messi", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Gojo Satoru", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Park Chan-wook", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
		{ID: utils.GetUniqueID(), CompanyID: companyID, Name: "Lee Myung-bak", Email: "", Description: "Generated Data", Address: "Jl. Kehidupan", PhoneNumber: "6281231231234", CreateBy: userID, UpdateBy: userID},
	}
	mapCustomers := []model.Customer{}
	for i, customer := range customers {
		mapCustomer := customer
		mapCustomer.Email = fmt.Sprintf("%s@gmail.com", strings.ToLower(strings.ReplaceAll(customer.Name, " ", "")))
		mapCustomer.Address = fmt.Sprintf("%s %d", mapCustomer.Address, i)
		mapCustomers = append(mapCustomers, mapCustomer)
	}
	tx.Create(&mapCustomers)

	tx = dbSeedOrderMajalah(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -1), 30)
	tx = dbSeedOrderMajalah(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -3), 20)
	tx = dbSeedOrderMajalah(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, 0), 50)
	tx = dbSeedOrderBatikModel(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -2), 3, 100)
	tx = dbSeedOrderBatikModel(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -1), 4, 50)
	tx = dbSeedOrderBatikModel(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -1), 6, 200)
	tx = dbSeedOrderBatikModel(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, 0), 2, 200)
	tx = dbSeedOrderBatikModel(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, 0), 5, 400)
	tx = dbSeedOrderKartuNama(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -3))
	tx = dbSeedOrderBanner(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -3))
	tx = dbSeedOrderNameTag(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -4))
	tx = dbSeedOrderStiker(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -4))
	tx = dbSeedOrderUndangan(tx, companyID, userID, mapCustomers[utils.GetRandomNumber(0, 100)%len(mapCustomers)], papers, phases, now.AddDate(0, 0, -5))

	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}

	fmt.Println("Seeding the database with initial data end")
}

var listTitleCompany = []string{
	"PLN",
	"BNI",
	"BRI",
	"BCA",
	"Pertamina",
	"Kejaksaan",
	"Kemenkeu",
	"Pegadaian",
	"Mandiri",
	"AEON",
	"Kemenham",
}

func dbSeedOrderMajalah(tx *gorm.DB, companyID, userID string, customer model.Customer, papers []model.Paper, phases []model.Phase, now time.Time, qty int64) *gorm.DB {
	title := listTitleCompany[utils.GetRandomNumber(0, 100)%len(listTitleCompany)]
	order := model.Order{
		ID:          utils.GetUniqueID(),
		CompanyID:   companyID,
		CustomerID:  customer.ID,
		Name:        fmt.Sprintf("Majalah %s 2 Model", title),
		Description: "",
		Number:      int64(utils.GetRandomNumber(100, 999)),
		CreateBy:    userID,
		CreateDt:    now,
		UpdateBy:    userID,
		UpdateDt:    now,
	}
	tx.Create(&order)

	prints := []model.Print{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			PaperID:     papers[5].ID,
			Name:        "Cover 1",
			Description: "Cover 1",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         qty,
			Price:       papers[5].DefaultPrice,
			Total:       1 * qty * papers[5].DefaultPrice,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			PaperID:     papers[6].ID,
			Name:        "Cover 2",
			Description: "Cover 2",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         qty,
			Price:       papers[6].DefaultPrice,
			Total:       1 * qty * papers[6].DefaultPrice,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			PaperID:     papers[3].ID,
			Name:        "Isi 1",
			Description: "Isi timbal balik 120 lembar",
			IsDuplex:    true,
			PageCount:   120,
			Qty:         qty,
			Price:       papers[3].DefaultPriceDuplex,
			Total:       120 * qty * papers[3].DefaultPriceDuplex,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			PaperID:     papers[3].ID,
			Name:        "Isi 2",
			Description: "Isi timbal balik 140 lembar",
			IsDuplex:    true,
			PageCount:   140,
			Qty:         qty,
			Price:       papers[3].DefaultPriceDuplex,
			Total:       140 * qty * papers[3].DefaultPriceDuplex,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
	}
	tx.Create(&prints)

	finishings := []model.Finishing{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Design Cover 1",
			Description: "Design Cover 1",
			Qty:         1,
			Price:       20000,
			Total:       1 * 20000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Design Cover 2",
			Description: "Design Cover 2",
			Qty:         1,
			Price:       30000,
			Total:       1 * 30000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Potong",
			Description: "Borongan potong 200K",
			Qty:         1,
			Price:       200000,
			Total:       1 * 200000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Laminating Cover",
			Description: "Laminating Cover 2 Model",
			Qty:         qty * 2,
			Price:       2400,
			Total:       qty * 2 * 2400,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Banding",
			Description: "Banding 2 Model",
			Qty:         qty * 2,
			Price:       17000,
			Total:       qty * 2 * 17000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
	}
	tx.Create(&finishings)

	transactions := []model.Transaction{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "DP",
			Description: "",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      qty * 400000,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 1),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 1),
		},
	}
	tx.Create(&transactions)

	orderphases := []model.Orderphase{
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[0].ID,
			Name:      phases[0].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 1),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 1),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[1].ID,
			Name:      phases[1].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 2),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 2),
		},
	}
	tx.Create(&orderphases)

	return tx
}

func dbSeedOrderBatikModel(tx *gorm.DB, companyID, userID string, customer model.Customer, papers []model.Paper, phases []model.Phase, now time.Time, modelCount, qty int64) *gorm.DB {
	title := listTitleCompany[utils.GetRandomNumber(0, 100)%len(listTitleCompany)]
	order := model.Order{
		ID:          utils.GetUniqueID(),
		CompanyID:   companyID,
		CustomerID:  customer.ID,
		Name:        fmt.Sprintf("Batik %s %d Model", title, modelCount),
		Description: fmt.Sprintf("Batik %s %d Model", title, modelCount),
		Number:      int64(utils.GetRandomNumber(100, 999)),
		CreateBy:    userID,
		CreateDt:    now,
		UpdateBy:    userID,
		UpdateDt:    now,
	}
	tx.Create(&order)

	prints := []model.Print{}
	for i := int64(1); i <= modelCount; i++ {
		prints = append(prints, model.Print{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			PaperID:     papers[6].ID,
			Name:        fmt.Sprintf("Model %s %d", title, i),
			Description: fmt.Sprintf("%d Lembar", qty),
			IsDuplex:    false,
			PageCount:   1,
			Qty:         qty,
			Price:       papers[6].DefaultPrice,
			Total:       1 * qty * papers[6].DefaultPrice,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		})
	}
	tx.Create(&prints)

	finishings := []model.Finishing{}
	tx.Create(&finishings)

	transactions := []model.Transaction{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Lunas",
			Description: "Lunas",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      modelCount * qty * papers[6].DefaultPrice,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 1),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 1),
		},
	}
	tx.Create(&transactions)

	orderphases := []model.Orderphase{
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[0].ID,
			Name:      phases[0].Name,
			CreateBy:  userID,
			CreateDt:  now.Add(1 * time.Hour),
			UpdateBy:  userID,
			UpdateDt:  now.Add(1 * time.Hour),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[1].ID,
			Name:      phases[1].Name,
			CreateBy:  userID,
			CreateDt:  now.Add(12 * time.Hour),
			UpdateBy:  userID,
			UpdateDt:  now.Add(12 * time.Hour),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[2].ID,
			Name:      phases[2].Name,
			CreateBy:  userID,
			CreateDt:  now.Add(24 * time.Hour),
			UpdateBy:  userID,
			UpdateDt:  now.Add(24 * time.Hour),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[5].ID,
			Name:      phases[5].Name,
			CreateBy:  userID,
			CreateDt:  now.Add(36 * time.Hour),
			UpdateBy:  userID,
			UpdateDt:  now.Add(36 * time.Hour),
		},
	}
	tx.Create(&orderphases)

	return tx
}

func dbSeedOrderKartuNama(tx *gorm.DB, companyID, userID string, customer model.Customer, papers []model.Paper, phases []model.Phase, now time.Time) *gorm.DB {
	order := model.Order{
		ID:          utils.GetUniqueID(),
		CompanyID:   companyID,
		CustomerID:  customer.ID,
		Name:        "Kartu Nama",
		Description: "4 box",
		Number:      int64(utils.GetRandomNumber(100, 999)),
		CreateBy:    userID,
		CreateDt:    now,
		UpdateBy:    userID,
		UpdateDt:    now,
	}
	tx.Create(&order)

	prints := []model.Print{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			PaperID:     papers[6].ID,
			Name:        "Kartu Nama",
			Description: "20 Lembar",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         20,
			Price:       papers[6].DefaultPrice,
			Total:       1 * 20 * papers[6].DefaultPrice,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
	}
	tx.Create(&prints)

	finishings := []model.Finishing{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     userID,
			Name:        "Potong",
			Description: "",
			Qty:         1,
			Price:       20000,
			Total:       1 * 20000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     userID,
			Name:        "Box  Kartu Nama",
			Description: "4 Box",
			Qty:         4,
			Price:       5000,
			Total:       4 * 5000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
	}
	tx.Create(&finishings)

	transactions := []model.Transaction{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "DP",
			Description: "",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      papers[6].DefaultPrice * 10,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 1),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 1),
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Lunas",
			Description: "",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      papers[6].DefaultPrice * 10,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 2),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 2),
		},
	}
	tx.Create(&transactions)

	orderphases := []model.Orderphase{
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[0].ID,
			Name:      phases[0].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 1),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 1),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[1].ID,
			Name:      phases[1].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 2),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 2),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[2].ID,
			Name:      phases[2].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 3),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 3),
		},
	}
	tx.Create(&orderphases)

	return tx
}

func dbSeedOrderBanner(tx *gorm.DB, companyID, userID string, customer model.Customer, papers []model.Paper, phases []model.Phase, now time.Time) *gorm.DB {
	order := model.Order{
		ID:          utils.GetUniqueID(),
		CompanyID:   companyID,
		CustomerID:  customer.ID,
		Name:        "Banner Ridwan Kamil",
		Description: "4 buah \n 4 x 12 ",
		Number:      int64(utils.GetRandomNumber(100, 999)),
		CreateBy:    userID,
		CreateDt:    now,
		UpdateBy:    userID,
		UpdateDt:    now,
	}
	tx.Create(&order)

	prints := []model.Print{}
	tx.Create(&prints)

	finishings := []model.Finishing{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Banner Ridwan Kamil",
			Description: "4 Box",
			Qty:         4,
			Price:       30000 * 4 * 12,
			Total:       4 * 30000 * 4 * 12,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
	}
	tx.Create(&finishings)

	transactions := []model.Transaction{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Lunas",
			Description: "",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      4 * 30000 * 4 * 12,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 2),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 2),
		},
	}
	tx.Create(&transactions)

	orderphases := []model.Orderphase{
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[0].ID,
			Name:      phases[0].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 1),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 1),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[4].ID,
			Name:      phases[4].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 2),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 2),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[5].ID,
			Name:      phases[5].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 3),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 3),
		},
	}
	tx.Create(&orderphases)

	return tx
}

func dbSeedOrderNameTag(tx *gorm.DB, companyID, userID string, customer model.Customer, papers []model.Paper, phases []model.Phase, now time.Time) *gorm.DB {
	order := model.Order{
		ID:          utils.GetUniqueID(),
		CompanyID:   companyID,
		CustomerID:  customer.ID,
		Name:        "Name Tag Peradi",
		Description: "10 buah",
		Number:      int64(utils.GetRandomNumber(100, 999)),
		CreateBy:    userID,
		CreateDt:    now,
		UpdateBy:    userID,
		UpdateDt:    now,
	}
	tx.Create(&order)

	prints := []model.Print{}
	tx.Create(&prints)

	finishings := []model.Finishing{
		{
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Design Foto",
			Description: "",
			Qty:         10,
			Price:       10000,
			Total:       10 * 10000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Name Tag Peradi",
			Description: "10 Buah",
			Qty:         10,
			Price:       20000,
			Total:       10 * 20000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
	}
	tx.Create(&finishings)

	transactions := []model.Transaction{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "DP",
			Description: "",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      10 * 10000,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 2),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 2),
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Lunas",
			Description: "",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      10 * 20000,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 2),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 2),
		},
	}
	tx.Create(&transactions)

	orderphases := []model.Orderphase{
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[0].ID,
			Name:      phases[0].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 1),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 1),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[4].ID,
			Name:      phases[4].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 2),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 2),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[5].ID,
			Name:      phases[5].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 3),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 3),
		},
	}
	tx.Create(&orderphases)

	return tx
}

func dbSeedOrderStiker(tx *gorm.DB, companyID, userID string, customer model.Customer, papers []model.Paper, phases []model.Phase, now time.Time) *gorm.DB {
	order := model.Order{
		ID:          utils.GetUniqueID(),
		CompanyID:   companyID,
		CustomerID:  customer.ID,
		Name:        "Stiker Selamat Makan",
		Description: "2 Model \n Jadi masing masing 300",
		Number:      int64(utils.GetRandomNumber(100, 999)),
		CreateBy:    userID,
		CreateDt:    now,
		UpdateBy:    userID,
		UpdateDt:    now,
	}
	tx.Create(&order)

	prints := []model.Print{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			PaperID:     papers[7].ID,
			Name:        "Model 1",
			Description: "Stiker Selamat Makan Model 1",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         30,
			Price:       papers[7].DefaultPrice,
			Total:       1 * 30 * papers[7].DefaultPrice,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			PaperID:     papers[7].ID,
			Name:        "Model 1",
			Description: "Stiker Selamat Makan Model 2",
			IsDuplex:    false,
			PageCount:   1,
			Qty:         30,
			Price:       papers[7].DefaultPrice,
			Total:       1 * 30 * papers[7].DefaultPrice,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
	}
	tx.Create(&prints)

	finishings := []model.Finishing{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Cutting",
			Description: "Borongan 200K",
			Qty:         1,
			Price:       200000,
			Total:       1 * 200000,
			CreateBy:    userID,
			CreateDt:    now,
			UpdateBy:    userID,
			UpdateDt:    now,
		},
	}
	tx.Create(&finishings)

	transactions := []model.Transaction{
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "DP",
			Description: "",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      papers[7].DefaultPrice * 30,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 1),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 1),
		},
		{
			ID:          utils.GetUniqueID(),
			CompanyID:   companyID,
			OrderID:     order.ID,
			Name:        "Lunas",
			Description: "",
			Type:        constant.TRANSACTION_TYPE_DEBIT,
			Amount:      papers[7].DefaultPrice*30 + 200000,
			CreateBy:    userID,
			CreateDt:    now.AddDate(0, 0, 2),
			UpdateBy:    userID,
			UpdateDt:    now.AddDate(0, 0, 2),
		},
	}
	tx.Create(&transactions)

	orderphases := []model.Orderphase{
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[0].ID,
			Name:      phases[0].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 1),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 1),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[2].ID,
			Name:      phases[2].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 2),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 2),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[3].ID,
			Name:      phases[3].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 3),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 3),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[4].ID,
			Name:      phases[4].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 4),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 4),
		},
		{
			ID:        utils.GetUniqueID(),
			CompanyID: companyID,
			OrderID:   order.ID,
			PhaseID:   phases[5].ID,
			Name:      phases[5].Name,
			CreateBy:  userID,
			CreateDt:  now.AddDate(0, 0, 5),
			UpdateBy:  userID,
			UpdateDt:  now.AddDate(0, 0, 5),
		},
	}
	tx.Create(&orderphases)

	return tx
}

func dbSeedOrderUndangan(tx *gorm.DB, companyID, userID string, customer model.Customer, papers []model.Paper, phases []model.Phase, now time.Time) *gorm.DB {
	return tx
}

func dbReset() {
	dbDown()
	dbUp()
	dbSeed()
}
