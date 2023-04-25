package data

var (
        //create statements
        putpacket = "insert into packet (pub_hash, created_at, voltage, frequency, protected, active_power, apparent_power, reactive_power, power_factor, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)"
        //select statements
        getpubconfigbyhash = "select pub_hash, nickname, kwp, kwpmake, kwr, kwrmake from pubconfig where pub_hash=$1 order by since desc limit 1"

        getpubbyhash = "select pub_id, created_at, latitude, longitude, hash from pub where hash=$1 order by created_at desc limit 1"

        getpubbyid = "select pub_id, created_at, latitude, longitude, hash from pub where pub_id=$1 order by created_at desc limit 1"

        getpubslimited = "select pub_id, created_at, latitude, longitude, hash from pub order by created_at desc limit $1"

        getpubsforsublimited = "select pub_id, created_at, latitude, longitude, hash from pub where creator=$1 order by created_at desc limit $2"

        getallpubsforsub = "select pub.pub_id, created_at, latitude, longitude, hash from pub inner join subpub on subpub.pub_id = pub.pub_id where sub_id=$1 order by created_at desc limit $2"

        getnumpubsforsub = "select count(*) from pub where creator=$1"

        getnumprotectedpubsforsub = "select count(*) from pub where creator=$1 and protected=true"

        getnumpubsonlineforsub = "select count(*) from pubconfig inner join pub on pub.hash = pubconfig.hash where pub.creator = $1 where pubconfig.since > $2"

        getpubfaultsforsub = "select pub_id, created_at, latitude, longitude, hash from pub where creator=$1 and protected=false order by created_at desc limit $2"

        getpubfaults = "select distinct on (pub_hash) packet.pub_hash, pubconfig.nickname, pubconfig.lastnotified from packet inner join pubconfig using(pub_hash) where pubconfig.notify=true and packet.protected=false order by pub_hash, created_at desc limit 10"

        getdummiesforsub = "select nickname, kwp, kwpmake, kwr, kwrmake, pub.creator, pub.latitude, pub.longitude from pubconfig inner join pub on pub.hash = pubconfig.pub_hash where pub.creator=$1 and pub.protected=false"

        getdummiesforall = "select nickname, kwp, kwpmake, kwr, kwrmake, pub.creator, pub.latitude, pub.longitude from pubconfig inner join pub on pub.hash = pubconfig.pub_hash where pub.protected=false"

        getsubforpub = "select sub.sub_id, sub.email, sub.name from sub inner join pub on sub.sub_id=pub.creator where pub.hash=$1 limit 1"

        getpubnameforpubhash = "select devicename from confo inner join pub on confo.hash = pub.hash where pub.hash=$1 limit 1"

        getsubdeetsbyemail = "select sub_id, created_at, email, name, phone from sub where email=$1 order by created_at desc limit 1"

        getsubpswdbyemail = "select sub_id, email, pswd from sub where email=$1 order by created_at desc limit 1"

        getsubslimited = "select sub_id, created_at, email, name, phone, verified from sub order by created_at desc limit $1"

        getcsubbyemail = "select sub_id, created_at, email from csub where email=$1 order by created_at desc limit 1"

        getverification = "select sub_id, email from sub where verification=$1"

        countsubs = "select count(distinct email) from sub"

        getlastconfobydevicename = "select devicename, ssid, created_at from confo where devicename=$1 and ssid=$2 order by created_at desc limit 1"

        getlastconfobyhash = "select devicename, ssid, created_at from confo where hash=$1 order by created_at desc limit 1"

        getpacketbyhash = "select id, created_at, voltage, frequency, protected from packet where pub_hash=$1 order by created_at desc limit 1"

        getlastpacketsbyhashlimited = "select id, created_at, voltage, frequency, protected, active_power, apparent_power, reactive_power, power_factor, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy from packet where pub_hash=$1 order by created_at desc limit $2"

        getpackets = "select pub_hash, created_at, voltage, frequency, protected, active_power, apparent_power, reactive_power, power_factor, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy from packet where created_at > $1 and created_at <= $2"

        getpacketsbyhash = "select pub_hash, created_at, voltage, frequency, protected, active_power, apparent_power, reactive_power, power_factor, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy from packet where pub_hash=$1 and created_at > $2 and created_at <= $3"

        gethourlies = "select pub_hash, timestamp, voltage_max, voltage_min, voltage_ave, voltage_exceptions, frequency_max, frequency_min, frequency_ave, frequency_exceptions, activepwr_max, activepwr_min, activepwr_ave, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy from hourly where timestamp > $1 and timestamp <= $2"

        getdailies = "select pub_hash, timestamp, voltage_max, voltage_min, voltage_ave, voltage_exceptions, frequency_max, frequency_min, frequency_ave, frequency_exceptions, activepwr_max, activepwr_min, activepwr_ave, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy from daily where timestamp > $1 and timestamp <= $2"

        gethourliesbyhash = "select pub_hash, timestamp, voltage_max, voltage_min, voltage_ave, voltage_exceptions, frequency_max, frequency_min, frequency_ave, frequency_exceptions, activepwr_max, activepwr_min, activepwr_ave, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy from hourly where timestamp > $1 and timestamp <= $2 and pub_hash=$3 order by timestamp asc"

        getdailiesbyhash = "select pub_hash, timestamp, voltage_max, voltage_min, voltage_ave, voltage_exceptions, frequency_max, frequency_min, frequency_ave, frequency_exceptions, activepwr_max, activepwr_min, activepwr_ave, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy from daily where timestamp > $1 and timestamp <= $2 and pub_hash=$3 order by timestamp asc"
        getlastsummarytime = "select pub_hash, timestamp from hourly order by timestamp desc limit 1"
        //update statements
        update_pub_protected_byhash = "update pub set protected=false where hash=$1"
        update_pubconfig_lastnotified_byhash = "update pubconfig set lastnotified=NOW() where pub_hash=$1"
        update_pubconfig_byhash = "update pubconfig set kwp = $1, kwpmake = $2, kwr = $3, kwrmake = $4, notify = $5 where pubconfig.pub_hash = $6"

        update_sub_pswd_byemail = "update sub set pswd=$1 where email=$2"
        update_sub_verified_byverification = "update sub set verified = TRUE where verification=$1"
        update_pub_byhash = "update pub set latitude = $1, longitude = $2, altitude = $3, orientation = $4, creator = $6 where pub.Hash = $5"
        //delete statements

)
