package data

import (
	//sq "database/sql"
	"fmt"
	_ "github.com/lib/pq"
        "github.com/golang/glog"
	//"golang.gurusys.co.uk/go-framework/sql"
        "crypto/sha1"
        "math/rand"
        "sort"
	"time"
)

var (
	dbcon *DB
        kwps = []float32{1.0, 3.0, 5.0, 10.0, 12.0, 20.0, 30.0, 50.0, 80.0, 100.0}
        kwpmakes = []string{"Canadian", "Emmvee", "First", "HHV", "JA", "Jinko", "Longi", "Q-Cells", "Waree", "Yingli"}
        kwrs = []float32{1.1, 3.3, 5.5, 10.8, 13.0, 21.0, 32.0, 53.0, 84.0, 105.0}
        kwrmakes = []string{"Delta", "Enphase", "Fronius", "Growatt", "Huawei", "K-Star", "SMA"}
        names = []string{"Absolutno", "Achird", "Acrab", "Adhara", "Adhil", "Albali", "Alderamin", "Algorab", "Alruba", "Atlas", "Bellatrix", "Capella", "Citadelle", "Copernicus", "Maia", "Markeb", "Mimosa", "Nahn", "Navi", "Nashira", "Nihal", "Polaris", "Bibha", "Revati", "Sarin", "Shaula", "Sirius", "Subra", "Vega", "Ashvini", "Bharani", "Kritika", "Rohini", "Mrigahirsha", "Ardra", "Punarvasu", "Pushya", "Ashlesha", "Magha", "Falguni", "Hasta", "Chitra", "Svati", "Vishakha", "Anuradha", "Jyeshtha", "Mula", "Ashadha", "Shravana", "Shatabhisha", "Dhanista", "Bhadrapada", "Revati", "Abhijit"}
)

func GetDB() (*DB, error) {
	var err error

	if dbcon != nil {
		return dbcon, nil
	}
	dbcon, err = Open()
	if err != nil {
		return nil, err
	}
	return dbcon, nil
}

type Coordinate struct {
        Latitude float32
        Longitude float32
        Altitude float32
}

type Confo struct {
        Id int64 `json:"id`
        Devicename string `json:"devicename"`
        Ssid string `json:"ssid"`
        Hash int64 `json:"hash"`
        //Created int64 `json:"created,omitempty"`
        Created time.Time `json:"created,omitempty"`
}

type Pub struct {
        Id int64 `json:"id"`
        Latitude float32 `json:"latitude,omitempty"`
        Longitude float32 `json:"longitude,omitempty"`
        Altitude float32 `json:"altitude,omitempty"`
        Orientation float32 `json:"orientation,omitempty"`
        Hash int64 `json:"hash"`
        Created time.Time `json:"created,omitempty"`
        Creator int64 `json:"email"`
        Protected bool `json:"protected"`
}

func (p *Pub) FormattedCreated() string {
        return p.Created.Format("2006-01-02 15:04:05")
}

type PubConfig struct {
        PubId int64 `json:"pubid"`
        Hash int64 `json:"hash"`
        Nickname string `json:"nickname,omitempty"`
        Typeref string `json:"typeref,omitempty"`
        Kwp float32 `json:"kwp,omitempty" form:"kwp"` // module
        Kwpmake string `json:"kwpmake,omitempty" form:"kwpm"`
        Kwr float32 `json:"kwr,omitempty" form:"kwr"` // inverter
        Kwrmake string `json:"kwrmake,omitempty" form:kwrm"`
        Kwlast float32 `json:"kwlast,omitempty"`
        Kwhhour float32 `json:"kwhhour,omitempty"`
        Kwhday float32 `json:"kwhday,omitempty"`
        Kwhlife float32 `json:"kwhlife,omitempty"`
        Since time.Time `json:"since,omitempty"`
        Visitslast time.Time `json:"visitslast,omitempty"`
        Visitslife int `json:"visitslife,omitempty"`
        LastUpdated time.Time
        Notify bool `json:"notify,omitempty" form:"notify"`
        LastNotified time.Time `json:"lastnotified,omitempty"`
}

func (p *PubConfig) FormattedLastNotified() string {
        return p.LastNotified.Format("2006-01-02 15:04:05")
}

type Dummies []*PubDummy
type PubDummy struct {
        Id int64 `json:"id"`
        Hash int64 `json:"hash"`
        Nickname string `json:"nickname,omitempty"`
        Latitude float32 `json:"latitude,omitempty"`
        Longitude float32 `json:"longitude,omitempty"`
        Created time.Time `json:"created,omitempty"`
        Kwp float32 `json:"kwp,omitempty"` // module
        Kwpmake string `json:"kwpmake,omitempty"`
        Kwr float32 `json:"kwr,omitempty"` // inverter
        Kwrmake string `json:"kwrmake,omitempty"`
        Kwlast float32 `json:"kwlast,omitempty"`
        Kwhday float32 `json:"kwhday,omitempty"`
        Kwhlife float32 `json:"kwhlife,omitempty"`
        Creator int64 `json:"creator,omitempty"`
        Email string `json:"email,omitempty"`
        Name string `json:"name,omitempty"`
        LastNotified time.Time `json:"lastnotified,omitempty"`
}
func (ds Dummies) Len() int {
        return len(ds)
}
func (ds Dummies) Swap(i, j int) {
        ds[i], ds[j] = ds[j], ds[i]
}
func (ds Dummies) Less(i,j int) bool {
        return ds[i].Kwlast > ds[j].Kwlast  // sorts descending
        //return ds[i].Kwlast < ds[j].Kwlast  // sorts asscending
}

type WrappedCoordinate struct {
        UserId int64
        Id int64
        Latitude float32
        Longitude float32
        Altitude float32
        Timestamp string
        Track string
}

type TrackRequest struct {
	User int64
	//Period *TimePeriod
	Track  string
}

// PutPub persists the provided Pub returning the pub_id
func PutPub(pub *Pub) (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        // convert to timestamp
        //created, err := time.Unix(confo.Created, 0).MarshalText()
        created, err := pub.Created.MarshalText()
	if err != nil {
                glog.Error(err)
		created, err = time.Now().MarshalText()
	}
        /*result, err := db.Exec("insert into pub (latitude, longitude, altitude, orientation, created_at, hash, creator) values ($1, $2, $3, $4, $5, $6, $7)", pub.Latitude, pub.Longitude, pub.Altitude, pub.Orientation, string(created), pub.Hash, pub.Creator)
        if err != nil {
                glog.Error(err)
                return 0 , err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Error("expected to affect 1 row, affected %d", rows)
                return uint64(rows) , err
        }
        return uint64(rows), nil*/
        result, err := db.Query("insert into pub (latitude, longitude, altitude, orientation, created_at, hash, creator) values ($1, $2, $3, $4, $5, $6, $7) returning pub_id", pub.Latitude, pub.Longitude, pub.Altitude, pub.Orientation, string(created), pub.Hash, pub.Creator)
        if err != nil {
                glog.Errorf("%v \n", err)
                return 0 , err
        }
	var id uint64
	defer result.Close()
	if !result.Next() {
                glog.Errorf("failed to insert any rows \n")
		return 0, fmt.Errorf("no rows returned on insert \n")
	}
	err = result.Scan(&id)
	if err != nil {
                fmt.Printf("failed to get id for new Pub:%s\n", err)
		return 0, fmt.Errorf("no id for new Pub (%s)", err)
	}
        return id, nil
}

// UpdatePub a Pub using hash of provided Pub
func UpdatePub(pub *Pub) error {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return err
        }
        //result, err := db.Exec("update pub set latitude = $1, longitude = $2, altitude = $3, orientation = $4, created_at = $5, creator = $7 where pub.Hash = $6", pub.Latitude, pub.Longitude, pub.Altitude, pub.Orientation, pub.Creator, pub.Hash, pub.Creator)
        result, err := db.Exec("update pub set latitude = $1, longitude = $2, altitude = $3, orientation = $4, creator = $6 where pub.Hash = $5", pub.Latitude, pub.Longitude, pub.Altitude, pub.Orientation, pub.Hash, pub.Creator)
        if err != nil {
                glog.Error("Couldn't update pub %v\n", err)
                return err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Error("Expected to affect 1 row, affected %d", rows)
                return err
        }
        return nil
}

func GetPubByHash(hash int64) (*Pub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select pub_id, created_at, latitude, longitude, hash from pub where hash=$1 order by created_at desc limit 1", hash)
        if err != nil {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, fmt.Errorf("No data for hash: %d \n", hash)
        }
        pc := &Pub{}
        err = rows.Scan(&pc.Id, &pc.Created, &pc.Latitude, &pc.Longitude, &pc.Hash)
        if err != nil {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, err
        }
        return pc, nil
}

func GetPubConfigByHash(hash int64) (*PubConfig, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select pub_hash, nickname, kwp, kwpmake, kwr, kwrmake, notify, lastnotified from pubconfig where pub_hash=$1 order by since desc limit 1", hash)
        if err != nil {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, fmt.Errorf("No data for hash: %d \n", hash)
        }
        pc := &PubConfig{}
        err = rows.Scan(&pc.Hash, &pc.Nickname, &pc.Kwp, &pc.Kwpmake, &pc.Kwr, &pc.Kwrmake, &pc.Notify, &pc.LastNotified)
        if err != nil {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, err
        }
        return pc, nil
}

func GetPubById(pub_id int64) (*Pub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select pub_id, created_at, latitude, longitude, hash from pub where pub_id=$1 order by created_at desc limit 1", pub_id)
        if err != nil {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, fmt.Errorf("No data for id: %d \n", pub_id)
        }
        pc := &Pub{}
        err = rows.Scan(&pc.Id, &pc.Created, &pc.Latitude, &pc.Longitude, &pc.Hash)
        if err != nil {
                glog.Errorf("data.GetPubByHash %v \n", err)
                return nil, err
        }
        return pc, nil
}

func GetPubs(limit int) ([]*Pub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select pub_id, created_at, hash, latitude, longitude, altitude, protected from pub order by created_at desc limit $1", limit)
        if err != nil {
                glog.Errorf("data.GetPubs %v \n", err)
                return nil, err
        }
        defer rows.Close()
        /*if !rows.Next() {
                glog.Errorf("data.GetPubs no rows \n")
                return nil, fmt.Errorf("No data for pub \n")
        }*/
        pbs := make([]*Pub, 0)
        for rows.Next() {
                pb := &Pub{}
                if err := rows.Scan(&pb.Id, &pb.Created, &pb.Hash, &pb.Latitude, &pb.Longitude, &pb.Altitude, &pb.Protected); err != nil {
                        glog.Errorf("data.GetPubs %v \n", err)
                        return pbs, fmt.Errorf("No data for pubs \n")
                }
                //glog.Infof("data.GetPubs appending \n")
                pbs = append(pbs, pb)
        }
        return pbs, nil
}

func GetPubDummies(limit int) (Dummies, error) {
        //pbs := make([]*PubDummy, 0)
        pbs := make(Dummies, 0)
        var lat,lng float32
        lat=13.0
        lng=77.5
        x := []float32{0.25, -0.25}
        //rand.Seed(time.Now().UnixNano())
        rand.Seed(123456)
        for i:=0; i<limit; i++ {
                la:= lat + rand.Float32() * x[rand.Intn(len(x))]
                lo:= lng + rand.Float32() * x[rand.Intn(len(x))]
                rh := time.Duration(-1 * rand.Intn(2400))
                ii := rand.Intn(len(kwps))
                kwp := kwps[ii]
                kwr := kwrs[ii]
                kw := kwp * 0.9
                now := time.Now().Add(time.Hour * rh).Round(time.Hour)
                pb := &PubDummy{Id: int64(i), Nickname: names[i], Latitude: la, Longitude: lo, Hash: int64(rand.Intn(10000)), Created: now, Kwp: kwp, Kwpmake: kwpmakes[rand.Intn(len(kwpmakes))], Kwr:kwr, Kwrmake: kwrmakes[rand.Intn(len(kwrmakes))], Kwlast: kw, Kwhday: kwp*4.5, Kwhlife: rand.Float32()*1000.0}
                pbs = append(pbs, pb)
        }
        sort.Sort(pbs)
        return pbs, nil
}

func GetPubsForSub(sub_id int64) ([]*Pub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select pub_id, created_at, latitude, longitude, hash, protected from pub where creator=$1 order by created_at desc limit $2", sub_id, 20)
        if err != nil {
                glog.Errorf("data.GetPubs %v \n", err)
                return nil, err
        }
        defer rows.Close()
        /*if !rows.Next() {
                glog.Errorf("data.GetPubs no rows \n")
                return nil, fmt.Errorf("No data for pub \n")
        }*/
        pbs := make([]*Pub, 0)
        for rows.Next() {
                pb := &Pub{}
                if err := rows.Scan(&pb.Id, &pb.Created, &pb.Latitude, &pb.Longitude, &pb.Hash, &pb.Protected); err != nil {
                        glog.Errorf("data.GetPubs %v \n", err)
                        return pbs, fmt.Errorf("No data for pubs \n")
                }
                //glog.Infof("data.GetPubs appending \n")
                pbs = append(pbs, pb)
        }
        return pbs, nil
}

//GetAllPubsForSub has a limit of 10 pubs in the query used
func GetAllPubsForSub(sub_id int64) ([]*Pub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select pub.pub_id, created_at, latitude, longitude, hash from pub inner join subpub on subpub.pub_id = pub.pub_id where sub_id=$1 order by created_at desc limit $2", sub_id, 10)
        if err != nil {
                glog.Errorf("data.GetPubs %v \n", err)
                return nil, err
        }
        defer rows.Close()
        /*if !rows.Next() {
                glog.Errorf("data.GetPubs no rows \n")
                return nil, fmt.Errorf("No data for pub \n")
        }*/
        pbs := make([]*Pub, 0)
        for rows.Next() {
                pb := &Pub{}
                if err := rows.Scan(&pb.Id, &pb.Created, &pb.Latitude, &pb.Longitude, &pb.Hash); err != nil {
                        glog.Errorf("data.GetPubs %v \n", err)
                        return pbs, fmt.Errorf("No data for pubs \n")
                }
                //glog.Infof("data.GetPubs appending \n")
                pbs = append(pbs, pb)
        }
        return pbs, nil
}

func GetPubFaultsForSub(sub_id int64) ([]*Pub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select pub_id, created_at, latitude, longitude, hash from pub where creator=$1 and protected=false order by created_at desc limit $2", sub_id, 20)
        if err != nil {
                glog.Errorf("data.GetPubs %v \n", err)
                return nil, err
        }
        defer rows.Close()
        /*if !rows.Next() {
                glog.Errorf("data.GetPubs no rows \n")
                return nil, fmt.Errorf("No data for pub \n")
        }*/
        pbs := make([]*Pub, 0)
        for rows.Next() {
                pb := &Pub{}
                if err := rows.Scan(&pb.Id, &pb.Created, &pb.Latitude, &pb.Longitude, &pb.Hash); err != nil {
                        glog.Errorf("data.GetPubs %v \n", err)
                        return pbs, fmt.Errorf("No data for pubs \n")
                }
                //glog.Infof("data.GetPubs appending \n")
                pbs = append(pbs, pb)
        }
        return pbs, nil
}

func GetUnilimitedPubFaultsForSub(sub_id int64) ([]*Pub, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select pub_id, created_at, latitude, longitude, hash from pub where creator=$1 and protected=false order by created_at desc", sub_id)
        if err != nil {
                glog.Errorf("data.GetPubs %v \n", err)
                return nil, err
        }
        defer rows.Close()
        /*if !rows.Next() {
                glog.Errorf("data.GetPubs no rows \n")
                return nil, fmt.Errorf("No data for pub \n")
        }*/
        pbs := make([]*Pub, 0)
        for rows.Next() {
                pb := &Pub{}
                if err := rows.Scan(&pb.Id, &pb.Created, &pb.Latitude, &pb.Longitude, &pb.Hash); err != nil {
                        glog.Errorf("data.GetPubs %v \n", err)
                        return pbs, fmt.Errorf("No data for pubs \n")
                }
                //glog.Infof("data.GetPubs appending \n")
                pbs = append(pbs, pb)
        }
        return pbs, nil
}

//GetDummiesForSub joins the Pub & PubConfig entries for a pub and creates dummies where protected = false
func GetDummiesForSub(sub_id int64) (Dummies, error) {
        pbds := make(Dummies, 0)
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return pbds, err
        }
        rows, err := db.Query(getdummiesforsub, sub_id)
        if err != nil {
                glog.Errorf("data.GetDummiesForSub dbquery %v \n", err)
                return nil, err
        }
        defer rows.Close()
        for rows.Next() {
                pb := &PubDummy{}
                if err := rows.Scan(&pb.Nickname, &pb.Kwp, &pb.Kwpmake, &pb.Kwr, &pb.Kwrmake, &pb.Latitude, &pb.Longitude); err != nil {
                        glog.Errorf("data.GetDummiesForSub %v \n", err)
                        return pbds, fmt.Errorf("No data for dummies \n")
                }
                pbds = append(pbds, pb)
        }
        return pbds, nil
}

func GetDummiesForAll() (Dummies, error) {
        pbds := make(Dummies, 0)
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return pbds, err
        }
        rows, err := db.Query(getdummiesforall)
        if err != nil {
                glog.Errorf("data.GetDummiesForSub dbquery %v \n", err)
                return nil, err
        }
        defer rows.Close()
        for rows.Next() {
                pb := &PubDummy{}
                if err := rows.Scan(&pb.Nickname, &pb.Kwp, &pb.Kwpmake, &pb.Kwr, &pb.Kwrmake, &pb.Creator, &pb.Latitude, &pb.Longitude); err != nil {
                        glog.Errorf("data.GetDummiesForSub %v \n", err)
                        return pbds, fmt.Errorf("No data for dummies \n")
                }
                pbds = append(pbds, pb)
        }
        return pbds, nil
}

//Populate fills the sub, pub & pubconfig tables with data required for testing/exhibiting
func Populate(subs, pubs int) ([]int64, error) {
        // count the subs
        s := CountSubs()
        // if less than desired populate subs and store a slice of sub ids
        sids := make([]int64, 0)
        hashes := make([]int64, 0)
        if s < subs {
               for i, s := range adopters {
                        _, err := PutSub(s) 
                        if err != nil {
                                glog.Errorf("putsub %d, %v \n", i, err)
                        } else {
                                ss, err := GetSubByEmail(s.Email)
                                if err != nil {
                                        glog.Errorf("getsubbyemail %d, %v \n", i, err)
                                } else {
                                        sids = append(sids, ss.Id)
                                }
                        }
                }
        }
        if len(sids) <= 0 {
                ss, err := GetSubs(10)
                if err != nil {
                        glog.Errorf("getsubs %v \n", err)
                }
                for _, s:=range ss {
                        sids = append(sids, s.Id)
                        glog.Infof("appending %v \n", len(sids))
                }
        }
        // generate random pubs distributed among the sub ids as creators - persist
        var lat,lng float32
        lat=13.0
        lng=77.5
        x := []float32{0.25, -0.25}
        //rand.Seed(time.Now().UnixNano())
        rand.Seed(123456)
        for i:=0; i<pubs; i++ {
                la:= lat + rand.Float32() * x[rand.Intn(len(x))]
                lo:= lng + rand.Float32() * x[rand.Intn(len(x))]
                rh := time.Duration(-1 * rand.Intn(2400))
                ii := rand.Intn(len(kwps))
                kwp := kwps[ii]
                kwr := kwrs[ii]
                kw := kwp * 0.9
                now := time.Now().Add(time.Hour * rh).Round(time.Hour)
                hash := int64(rand.Intn(10000))
                hashes = append(hashes, hash)
                pb := &Pub{Latitude: la, Longitude: lo, Hash: hash, Creator: sids[rand.Intn(len(sids))], Created: now, Protected: true}
                _, err := PutPub(pb)
                if err != nil {
                        glog.Errorf("populate putpub %v \n", err)
                } else {
                        // generate a confo for the pub produced in previous - persist
                        pc := &Confo{Devicename:names[i], Ssid:names[i+1], Hash:hash}
                        _, err := PutConfo(pc)
                        if err != nil {
                                glog.Errorf("populate putconfo %v \n", err)
                        }
                        // generate a pubconfig for the pub produced in previous - persist
                        pbc := &PubConfig{ Nickname: names[i], Hash: hash, Kwp: kwp, Kwpmake: kwpmakes[rand.Intn(len(kwpmakes))], Kwr:kwr, Kwrmake: kwrmakes[rand.Intn(len(kwrmakes))], Kwlast: kw, Kwhday: kwp*4.5, Kwhlife: rand.Float32()*1000.0}
                        _, err = PutPubConfig(pbc)
                        if err != nil {
                                glog.Errorf("populate putpubc %v \n", err)
                        }
                }
        }
        return hashes, nil
}

type PubStat struct {
        T int64
        O int64
        P int64
        P1 float64
        P2 float64
        P3 float64
        P4 float64
}

//GetPubStatsForSub queries db for count of total, protected pubs by sub and last power readings of all pubs by sub - TODO - there should be a separate single table getting updated with this data in the background such that this user query hits only one table once. 
func GetPubStatsForSub(sub_id int64) (*PubStat, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query(getnumpubsforsub, sub_id)
        if err != nil {
                glog.Errorf("data.GetPubStatsForSub dbquery getnumpubsforsub %v \n", err)
                return nil, err
        }
        defer rows.Close()
        ps := &PubStat{}
        rows.Next() // should be only one row
        err = rows.Scan(&ps.T)
        if err != nil {
                glog.Errorf("data.GetPubStatsForSub rowscan %v \n", err)
                return nil, err
        }

        //ps.O=ps.T
        rows, err = db.Query(getnumpubsonlineforsub, sub_id, time.Now().Add(time.Hour*-24))
        if err != nil {
                glog.Errorf("data.GetPubStatsForSub dbquery getnumpubsonlineforsub %v \n", err)
                return nil, err
        }
        defer rows.Close()
        rows.Next() // should be only one row
        err = rows.Scan(&ps.O)
        if err != nil {
                glog.Errorf("data.GetPubStatsForSub rowscan getnumpubsonlineforsub %v \n", err)
                return nil, err
        }

        rows, err = db.Query(getnumprotectedpubsforsub, sub_id)
        if err != nil {
                glog.Errorf("data.GetPubStatsForSub dbquery getnumprotectedpubsforsub %v \n", err)
                return nil, err
        }
        defer rows.Close()
        rows.Next() // should be only one row
        err = rows.Scan(&ps.P)
        if err != nil {
                glog.Errorf("data.GetPubStatsForSub rowscan getnumprotectedpubsforsub %v \n", err)
                return nil, err
        }
        return ps, nil
}

// GetPubFaults queries table packet for the latest `unprotected` packet from a pub in the latest 100 packets received where notify config for that pub is true. Either fix, turn notify off or receive daily emails. 
func GetPubFaults(withCreator bool) (Dummies, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select distinct on (pub_hash) packet.pub_hash, pubconfig.nickname, pubconfig.lastnotified from packet inner join pubconfig using(pub_hash) where pubconfig.notify=true and packet.protected=false order by pub_hash, created_at desc limit 10")
        if err != nil {
                glog.Errorf("data.GetPubs %v \n", err)
                return nil, err
        }
        defer rows.Close()
        /*if !rows.Next() {
                glog.Errorf("data.GetPubs no rows \n")
                return nil, fmt.Errorf("No data for pub \n")
        }*/
        pbs := make(Dummies, 0)
        for rows.Next() {
                pb := &PubDummy{}
                if err := rows.Scan(&pb.Hash, &pb.Nickname, &pb.LastNotified); err != nil {
                        glog.Errorf("data.GetPubFaults %v \n", err)
                        return pbs, fmt.Errorf("No data for pubfaults \n")
                }
                // append only if last notified more than 24 hours ago
                if pb.LastNotified.Before(time.Now().Add(time.Hour * -24)) {
                        pbs = append(pbs, pb)
                }
        }
        if len(pbs) == 0 {
                glog.Infof("No faulty packets found! ")
        }
        // still need to populate the creator, email, name fields
        if withCreator {
                for _, pb:= range pbs {
                        rows1, err := db.Query("select sub.sub_id, sub.email, sub.name from sub inner join pub on sub.sub_id=pub.creator where pub.hash=$1 limit 1", pb.Hash)
                        if err != nil {
                                glog.Errorf("data.GetPubFaults wcreator %v \n", err)
                                return nil, err
                        }
                        defer rows1.Close()
                        for rows1.Next() {
                                err = rows1.Scan(&pb.Creator, &pb.Email, &pb.Name)
                                if err != nil {
                                        glog.Errorf("data.GetPubFaults wcreator scan %v \n", err)
                                        return nil, err
                                }
                        }
                }
        }
        return pbs, nil
}

//UpdatePubStatus updates `pub.protected` and `pubconfig.lastnotified` in two separate update queries
func UpdatePubStatus(pub_hash int64, status, setlastnotified bool) error {
        db, err := GetDB()
        if err != nil {
                glog.Errorf("updatepubstatus getdb %v \n", err)
                return err
        }
        result, err := db.Exec("update pub set protected=$1 where hash=$2", status, pub_hash)
        if err != nil {
                glog.Errorf("updatebpubstatus update %v \n", err)
                return err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Errorf("Expected to affect 1 row, affected %d", rows)
                return fmt.Errorf("updatepubstatus no rows updated")
        }
        if setlastnotified {
                result, err = db.Exec("update pubconfig set lastnotified=NOW() where pub_hash=$1", pub_hash)
                if err != nil {
                        glog.Errorf("updatebpubstatus setlastnotified %v \n", err)
                        return err
                }
                rows, err = result.RowsAffected()
                if rows != 1 {
                        glog.Errorf("Expected to affect 1 row, affected %d", rows)
                        return fmt.Errorf("updatepubstatus setlastnotified no rows updated")
                }
        } else {
                result, err = db.Exec("update pubconfig set since=NOW() where pub_hash=$1", pub_hash)
                if err != nil {
                        glog.Errorf("updatebpubstatus setsince %v \n", err)
                        return err
                }
                rows, err = result.RowsAffected()
                if rows != 1 {
                        glog.Errorf("Expected to affect 1 row, affected %d", rows)
                        return fmt.Errorf("updatepubstatus setsince no rows updated")
                }
        }
        return nil
}

// GetPubDeviceName returns DeviceName for `pub_hash`from `confo`
func GetPubDeviceName(pub_hash int64) (string, error){
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return "", err
        }
        rows, err := db.Query("select devicename from confo inner join pub on confo.hash = pub.hash where pub.hash=$1 limit 1", pub_hash)
        if err != nil {
                glog.Errorf("data.GetPubDeviceName %v \n", err)
                return "", err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetPubDeviceName %v \n", err)
                return "", fmt.Errorf("No devicename for: %d \n", pub_hash)
        }
        devicename := ""
        err = rows.Scan(&devicename)
        if err != nil {
                glog.Errorf("data.GetPubDeviceName %v \n", err)
                return "", err
        }
        return devicename, nil;
}

// PutPubForSub populates the 'pubsub' table with the supplied sub_id and pub_id
func PutPubForSub(sub_id int, pub_id int) (int, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        result, err := db.Exec("insert into subpub (sub_id, pub_id) values ($1, $2)", sub_id, pub_id)
        if err != nil {
                glog.Error(err)
                return 0 , err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Error("expected to affect 1 row, affected %d", rows)
                return int(rows) , err
        }
        return int(rows), nil
}

// PutPubConfig puts the provided PubConfig into pg table `pubconfig`
func PutPubConfig(pubc *PubConfig) (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        result, err := db.Exec("insert into pubconfig (pub_hash, nickname, kwp, kwpmake, kwr, kwrmake, notify) values ($1, $2, $3, $4, $5, $6, $7)", pubc.Hash, pubc.Nickname, pubc.Kwp, pubc.Kwpmake, pubc.Kwr, pubc.Kwrmake, pubc.Notify)
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

// UpdatePubConfig updates a PubConfig in table `pubconfig` using hash of provided PubConfig
func UpdatePubConfig(pubc *PubConfig) error {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return err
        }
        if pubc.Hash == 0 {
                glog.Error("update pubconfig no hash provided  \n")
                return fmt.Errorf("invald hash : %d \n", pubc.Hash)
        }
        result, err := db.Exec("update pubconfig set kwp = $1, kwpmake = $2, kwr = $3, kwrmake = $4, notify = $5 where pubconfig.pub_hash = $6", pubc.Kwp, pubc.Kwpmake, pubc.Kwr, pubc.Kwrmake, pubc.Notify, pubc.Hash)
        if err != nil {
                glog.Errorf("Couldn't update pub %v \n", err)
                return err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Errorf("Expected to affect 1 row, affected %d", rows)
                if rows == 0 {
                        if pubc.Hash == 0 {
                                glog.Errorf("Couln't create pubconfig %d hash", pubc.Hash)
                                return fmt.Errorf("Couln't create pubconfig %d hash", pubc.Hash)
                        }
                        nick, err := GetPubDeviceName(pubc.Hash)
                        if err != nil {
                                glog.Errorf("Couln't get pub name of %d hash", pubc.Hash)
                                return fmt.Errorf("Couln't create pubconfig %d hash", pubc.Hash)
                        }
                        pubc.Nickname = nick
                        if _, err = PutPubConfig(pubc); err != nil {
                                glog.Errorf("Couln't put pubconfig %v", err)
                                return err
                        }
                }
                return err
        }
        return nil
}

// PutConf inserts a recd. Conf in db. 
func PutConfo(confo *Confo) (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        // check and convert to timestamp
        //created, err := time.Unix(confo.Created, 0).MarshalText()
        created, err := confo.Created.MarshalText()
	if err != nil || confo.Created.Before(time.Date(2006,1,2,3,4,5,0,time.UTC)) {
                if err != nil {
                        glog.Error(err)
                }
		created, err = time.Now().MarshalText()
	}
        result, err := db.Exec("insert into confo (devicename, ssid, created_at, hash) values ($1, $2, $3, $4)", confo.Devicename, confo.Ssid, string(created), confo.Hash)
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

// GetLastConf retrieves the latest conf in db with supplied `ssid+devicename`. 
// It returns a *Conf with matching `ssid+devicename` and latest timestamp or nil, error if none found
func GetLastConfo(devicename, ssid string) (*Confo, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select devicename, ssid, created_at, hash from confo where devicename=$1 and ssid=$2 order by created_at desc limit 1", devicename, ssid)
        if err != nil {
                glog.Errorf("data.GetLastConfo %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetLastConfo %v \n", err)
                return nil, fmt.Errorf("No data for tuple: %s %s \n", devicename, ssid)
        }
        pc := &Confo{}
        err = rows.Scan(&pc.Devicename, &pc.Ssid, &pc.Created, &pc.Hash)
        if err != nil {
                glog.Errorf("data.GetLastConfo %v \n", err)
                return nil, err
        }
        return pc, nil
}

func GetLastConfoWithHash(hash int64) (*Confo, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select devicename, ssid, created_at from confo where hash=$1 order by created_at desc limit 1", hash)
        if err != nil {
                glog.Errorf("data.GetLastConfoWithHash %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetLastConfoWithHash %v \n", err)
                return nil, fmt.Errorf("No data for : %d \n", hash)
        }
        pc := &Confo{}
        err = rows.Scan(&pc.Devicename, &pc.Ssid, &pc.Created)
        if err != nil {
                glog.Errorf("data.GetLastConfoWithHash %v \n", err)
                return nil, err
        }
        return pc, nil
}

// PutPacket inserts packet in db. It *should* check whether pub_id sent with packet is from the 'correct' pub. This could be a check against the location or a 'secret' decided during configuration with app which is sent along with each packet'
func PutPacket(packet *Packet) (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        result, err := db.Exec("insert into packet (pub_hash, created_at, voltage, frequency, protected, active_power, apparent_power, reactive_power, power_factor, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", packet.Id, packet.Timestamp, packet.Voltage, packet.Frequency, packet.Status, packet.ActiPwr, packet.AppaPwr, packet.ReacPwr, packet.PwrFctr, packet.ImActEn, packet.ExActEn, packet.ImRctEn, packet.ExRctEn, packet.TlActEn, packet.TlRctEn)
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

// GetPacket retrieves the latest packet in db with supplied `pub_hash`. 
// Note: it returns a packet with `id` set to the serial of the reading rather than `pub_hash` since caller of function already has `pub_hash`
func GetLastPacket(pubHash int64) (*Packet, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select id, created_at, voltage, frequency, protected from packet where pub_hash=$1 order by created_at desc limit 1", pubHash)
        if err != nil {
                glog.Errorf("data.GetLastPacket %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetLastPacket %v \n", err)
                return nil, fmt.Errorf("No data for user: %d \n", pubHash)
        }
        pc := &Packet{Id: pubHash}
        err = rows.Scan(&pc.Id, &pc.Timestamp, &pc.Voltage, &pc.Frequency, &pc.Status)
        if err != nil {
                glog.Errorf("data.GetLastPacket %v \n", err)
                return nil, err
        }
        return pc, nil
}

func GetLastPackets(pubHash int64, limit int) ([]*Packet, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select id, created_at, voltage, frequency, protected, active_power, apparent_power, reactive_power, power_factor, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy from packet where pub_hash=$1 order by created_at desc limit $2", pubHash, limit)
        if err != nil {
                glog.Errorf("data.GetLastPacket %v \n", err)
                return nil, err
        }
        defer rows.Close()
        pcks := make([]*Packet, 0)
        for rows.Next() {
                pck := &Packet{}
                if err := rows.Scan(&pck.Id, &pck.Timestamp, &pck.Voltage, &pck.Frequency, &pck.Status, &pck.ActiPwr, &pck.AppaPwr, &pck.ReacPwr, &pck.PwrFctr, &pck.ImActEn, &pck.ExActEn, &pck.ImRctEn, &pck.ExRctEn, &pck.TlActEn, &pck.TlRctEn); err != nil {
                        glog.Errorf("data.GetLastPackets %v \n", err)
                        return pcks, fmt.Errorf("No data for packets \n")
                }
                //glog.Infof("data.GetSubs appending \n")
                pcks = append(pcks, pck)
        }
        if len(pcks) == 0 {
                pcks = GetDummyPackets(10)
        }
        return pcks, nil
}

//GetPackets queries for packets for all pubs from (exclusive) to (inclusive) times provided as argument
func GetPackets(from, to time.Time) ([]*Packet, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query(getpackets, from, to)
        if err != nil {
                glog.Errorf("data.GetLastPacket %v \n", err)
                return nil, err
        }
        defer rows.Close()
        pcks := make([]*Packet, 0)
        for rows.Next() {
                pck := &Packet{}
                if err := rows.Scan(&pck.Id, &pck.Timestamp, &pck.Voltage, &pck.Frequency, &pck.Status, &pck.ActiPwr, &pck.AppaPwr, &pck.ReacPwr, &pck.PwrFctr, &pck.ImActEn, &pck.ExActEn, &pck.ImRctEn, &pck.ExRctEn, &pck.TlActEn, &pck.TlRctEn); err != nil {
                        glog.Errorf("data.GetLastPackets %v \n", err)
                        return pcks, fmt.Errorf("No data for packets \n")
                }
                //glog.Infof("data.GetSubs appending \n")
                pcks = append(pcks, pck)
        }
        return pcks, nil
}

//GetPacketsByHash queries for packets from pub with pubHash, from (exclusive) to (inclusive) times provided as argument
func GetPacketsByHash(pubHash int64, from, to time.Time) ([]*Packet, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query(getpacketsbyhash, pubHash, from, to)
        if err != nil {
                glog.Errorf("data.GetLastPacket %v \n", err)
                return nil, err
        }
        defer rows.Close()
        pcks := make([]*Packet, 0)
        for rows.Next() {
                pck := &Packet{}
                if err := rows.Scan(&pck.Id, &pck.Timestamp, &pck.Voltage, &pck.Frequency, &pck.Status, &pck.ActiPwr, &pck.AppaPwr, &pck.ReacPwr, &pck.PwrFctr, &pck.ImActEn, &pck.ExActEn, &pck.ImRctEn, &pck.ExRctEn, &pck.TlActEn, &pck.TlRctEn); err != nil {
                        glog.Errorf("data.GetLastPackets %v \n", err)
                        return pcks, fmt.Errorf("No data for packets \n")
                }
                //glog.Infof("data.GetSubs appending \n")
                pcks = append(pcks, pck)
        }
        return pcks, nil
}

// PutCoordinate
func PutCoordinate(coord *WrappedCoordinate) (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        //result, err := db.Exec("insert into coordinate (user_id, latitude, longitude, altitude, created_at) values ($1, $2, $3, $4, $5)", coord.UserId, coord.Latitude, coord.Longitude, coord.Altitude, coord.Timestamp)
        result, err := db.Exec("insert into coordinate (user_id, latitude, longitude, altitude, track) values ($1, $2, $3, $4, $5)", coord.UserId, coord.Latitude, coord.Longitude, coord.Altitude, coord.Track)
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

func GetCoordinate(userid int64) (*WrappedCoordinate, error){
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        rows, err := db.Query("select id, latitude, longitude, altitude from coordinate where user_id=$1 order by created_at desc limit 1", userid)
        if err != nil {
                glog.Errorf("data.GetCoordinate %v \n", err)
                return nil, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetCoordinate %v \n", err)
                return nil, fmt.Errorf("No data for user: %d \n", userid)
        }
        wc := &WrappedCoordinate{UserId: userid}
        err = rows.Scan(&wc.Id, &wc.Latitude, &wc.Longitude, &wc.Altitude)
        if err != nil {
                glog.Errorf("data.GetCoordinate %v \n", err)
                return nil, err
        }
        return wc, nil
}

func GetTrack(tr *TrackRequest) ([]*Coordinate, error){
        cs := make([]*Coordinate, 0)
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return cs, err
        }
        if tr.User == 0 || tr.Track == ""{
                return nil, fmt.Errorf("Invalid user or track")
        }
        rows, err := db.Query("select latitude, longitude, altitude from coordinate where user_id=$1and track=$2 order by created_at asc", tr.User, tr.Track)
        if err != nil {
                glog.Errorf("data.GetCoordinate %v \n", err)
                return cs, err
        }
        defer rows.Close()
        if !rows.Next() {
                glog.Errorf("data.GetCoordinate %v \n", err)
                return cs, fmt.Errorf("No data for user, track: %d %s \n", tr.User, tr.Track)
        }
        for rows.Next() {
                c := &Coordinate{}
                if err := rows.Scan(&c.Latitude, &c.Longitude, &c.Altitude); err != nil {
                        glog.Errorf("data.GetCoordinate %v \n", err)
                        return cs, fmt.Errorf("No data for user: %d \n", tr.User)
                }
                cs = append(cs, c)
        }
        return cs, nil
}

// Sha1Str returns the first 8 bytes in hex string representation of the SHA1 hash of the input str
func Sha1Str(str string) string {
        eb := []byte(str)
        ebs20 := sha1.Sum(eb)
        ebs8 := ebs20[:8]
        sha1str := fmt.Sprintf("%x", ebs8)
        return sha1str
}
