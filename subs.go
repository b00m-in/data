package data

import (
        "fmt"
        "time"
        "github.com/golang/glog"
        "golang.org/x/crypto/bcrypt"
)


var (
        emails = []string{"rcs@pholtovolta.in","user@gr00v.in","coo@l0v.in",""}
        snames = []string{"Carlito","Roger","Rafa","Novak"}
        adopters = []*Sub{
                &Sub{Email:emails[0], Phone:"23456", Name: snames[0], Pswd:"asdf", Created: time.Now(), Verification: "asdf1234", Verified: true},
                &Sub{Email:emails[1], Phone:"23456", Name: snames[1], Pswd:"asdf", Created: time.Now(), Verification: "dfgh1234", Verified: true},
                &Sub{Email:emails[2], Phone:"23456", Name: snames[2], Pswd:"asdf", Created: time.Now(), Verification: "ddfg1234", Verified: true},
                }
)

type Sub struct {
        Id int64 `json:"id"`
        Email string `json:"email" form:"email"`
        Name string `json:"name" form:"name"`
        Phone string `json:"phone" form:"phone"`
        Pswd string `json:"pswd" form:"pswd"`
        Created time.Time `json:"created,omitempty"`
        Verification string `json:"verification,omitempty"`
        Verified bool `json:"verified"`
}

func (s *Sub) FormattedCreated() string {
        return s.Created.Format("2006-01-02 15:04:05")
}

func (sub *Sub) Put() (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        // convert to timestamp
        //created, err := time.Unix(confo.Created, 0).MarshalText()
        created, err := sub.Created.MarshalText()
	if err != nil || sub.Created.Before(time.Date(2000,1,1,1,1,1,1,time.UTC)) {
                glog.Error(err)
		created, err = time.Now().MarshalText()
	}
        pb, err := Hash([]byte(sub.Pswd), bcrypt.DefaultCost)
        if err != nil {
                glog.Errorf("%v\n", err)
        } else {
                sub.Pswd = string(pb)
        }
        result, err := db.Exec("insert into sub (email, phone, name, pswd, created_at, verification) values ($1, $2, $3, $4, $5, $6)", sub.Email, sub.Phone, sub.Name, sub.Pswd, string(created), sub.Verification)
        if err != nil {
                glog.Error(err)
                return 0 , err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Error("expected to affect 1 row, affected %d", rows)
                return uint64(rows) , err
        }
        return uint64(rows), nil
}

func (sub *Sub) Update() error {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return err
        }
        pb, err := Hash([]byte(sub.Pswd), bcrypt.DefaultCost)
        if err != nil {
                glog.Errorf("%v\n", err)
        } else {
                sub.Pswd = string(pb)
        }
        result, err := db.Exec("update sub set pswd=$1 where email=$2", sub.Pswd, sub.Email)
        if err != nil {
                glog.Errorf("Couldn't update pub %v \n", err)
                return err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Errorf("Expected to affect 1 row, affected %d", rows)
                return fmt.Errorf("sub %s not updated", sub.Email)
        }
        return nil
}

func CountSubs() int {
        db, err := GetDB()
        if err != nil {
                glog.Errorf("countsubs getdb %v \n", err)
                return 0
        }

        rows, err := db.Query(countsubs)
        if err != nil {
                glog.Errorf("countsubs query %v \n", err)
                return 0
        }
        defer rows.Close()
        var cs int
        rows.Next() // should be only one row
        err = rows.Scan(&cs)
        if err != nil {
                glog.Errorf("countsubs rowscan %v \n", err)
                return 0
        }
        return cs
}

func Hash(pswd []byte, cost int) ([]byte, error) {
        hash, err := bcrypt.GenerateFromPassword(pswd, bcrypt.DefaultCost)
        if err != nil {
                return nil, err
        }
        return hash, nil
}

func CompareHash(email string, pswd string) bool {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return false
        }
        rows, err := db.Query("select sub_id, email, pswd from sub where email=$1 order by created_at desc limit 1", email)
        if err != nil {
                glog.Errorf("data.CheckPswd %v \n", err)
                return false
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.CheckPswd %v \n", err)
                return false
        }
        pc := &Sub{}
        err = rows.Scan(&pc.Id, &pc.Email, &pc.Pswd)
        if err != nil {
                glog.Errorf("data.CheckPswd %v \n", err)
                return false
        }
        /*if pc.Pswd != pswd {
                return false
        }*/
        if err := bcrypt.CompareHashAndPassword([]byte(pc.Pswd), []byte(pswd)); err != nil {
                glog.Errorf("data.CheckPswd %v \n", err)
                return false
        }
        return true
}

func PutSub(sub *Sub) (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        // convert to timestamp
        //created, err := time.Unix(confo.Created, 0).MarshalText()
        created, err := sub.Created.MarshalText()
	if err != nil || sub.Created.Before(time.Date(2000,1,1,1,1,1,1,time.UTC)) {
                glog.Error(err)
		created, err = time.Now().MarshalText()
	}
        result, err := db.Exec("insert into sub (email, phone, name, pswd, created_at, verification) values ($1, $2, $3, $4, $5, $6)", sub.Email, sub.Phone, sub.Name, sub.Pswd, string(created), sub.Verification)
        if err != nil {
                glog.Error(err)
                return 0 , err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Error("expected to affect 1 row, affected %d", rows)
                return uint64(rows) , err
        }
        return uint64(rows), nil
}

func UpdateSub(sub *Sub) error {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return err
        }
        result, err := db.Exec("update sub set pswd=$1 where email=$2", sub.Pswd, sub.Email)
        if err != nil {
                glog.Errorf("Couldn't update pub %v \n", err)
                return err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Errorf("Expected to affect 1 row, affected %d", rows)
                return fmt.Errorf("sub %s not updated", sub.Email)
        }
        return nil
}

func GetSubByEmail(email string) (*Sub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select sub_id, created_at, email, name, phone from sub where email=$1 order by created_at desc limit 1", email)
        if err != nil {
                glog.Errorf("data.GetSubByEmail %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetSubByEmail %v \n", err)
                return nil, fmt.Errorf("No data for email: %s \n", email)
        }
        pc := &Sub{}
        err = rows.Scan(&pc.Id, &pc.Created, &pc.Email, &pc.Name, &pc.Phone)
        if err != nil {
                glog.Errorf("data.GetSubByEmail %v \n", err)
                return nil, err
        }
        return pc, nil
}

func CheckPswd(email string, pswd string) bool {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return false
        }
        rows, err := db.Query("select sub_id, email, pswd from sub where email=$1 order by created_at desc limit 1", email)
        if err != nil {
                glog.Errorf("data.CheckPswd %v \n", err)
                return false
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.CheckPswd %v \n", err)
                return false
        }
        pc := &Sub{}
        err = rows.Scan(&pc.Id, &pc.Email, &pc.Pswd)
        if err != nil {
                glog.Errorf("data.CheckPswd %v \n", err)
                return false
        }
        if pc.Pswd != pswd {
                return false
        }
        return true
}

func GetSubs(limit int) ([]*Sub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select sub_id, created_at, email, name, phone, verified from sub order by created_at desc limit $1", limit)
        if err != nil {
                glog.Errorf("data.GetSubs %v \n", err)
                return nil, err
        }
        defer rows.Close()
        /*if !rows.Next() {
                glog.Errorf("data.GetPubs no rows \n")
                return nil, fmt.Errorf("No data for pub \n")
        }*/
        sbs := make([]*Sub, 0)
        for rows.Next() {
                sb := &Sub{}
                if err := rows.Scan(&sb.Id, &sb.Created, &sb.Email, &sb.Name, &sb.Phone, &sb.Verified); err != nil {
                        glog.Errorf("data.GetSubs %v \n", err)
                        return sbs, fmt.Errorf("No data for subs \n")
                }
                //glog.Infof("data.GetSubs appending \n")
                sbs = append(sbs, sb)
        }
        return sbs, nil
}

// PutCsub persists an unknown Sub with unregistered email which may be part of a confo from device
func PutCsub(sub *Sub) (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        result, err := db.Exec("insert into csub (email) values ($1)", sub.Email)
        if err != nil {
                glog.Error(err)
                return 0 , err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Error("expected to affect 1 row, affected %d", rows)
                return uint64(rows) , err
        }
        return uint64(rows), nil
}

func GetCsubByEmail(email string) (*Sub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select sub_id, created_at, email from csub where email=$1 order by created_at desc limit 1", email)
        if err != nil {
                glog.Errorf("data.GetCsubByEmail %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetCsubByEmail %v \n", err)
                return nil, fmt.Errorf("No data for email: %s \n", email)
        }
        pc := &Sub{}
        err = rows.Scan(&pc.Id, &pc.Created, &pc.Email)
        if err != nil {
                glog.Errorf("data.GetCsubByEmail %v \n", err)
                return nil, err
        }
        return pc, nil
}

func CheckVerification(verification string) (*Sub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select sub_id, email from sub where verification=$1", verification)
        if err != nil {
                glog.Errorf("data.CheckVerification %v \n", err)
                return nil, err
        }
        if !rows.Next() {
                glog.Errorf("data.CheckVerification %v \n", err)
                return nil, fmt.Errorf("No data for verification: %s \n", verification)
        }
        pc := &Sub{}
        err = rows.Scan(&pc.Id, &pc.Email)
        if err != nil {
                glog.Errorf("data.CheckVerification %v \n", err)
                return nil, err
        }
        rows.Close()
        _, err = db.Exec("update sub set verified = TRUE where verification=$1", verification)
        if err != nil {
                glog.Errorf("data.CheckVerification %v \n", err)
                return nil, err
        }
        return pc, nil
}
