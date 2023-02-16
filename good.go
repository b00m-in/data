package data

import (
        //"fmt"

)

type Category struct {
        Name string `json:"category"`
        Subcategories []string `json:"subcategories"`
}
type Good struct {
        Id int64 `json:"id,string" schema:"-"` //1,000-9,999
        Code string `json:"code" schema:"code"`
        Category string `json:"category" schema:"category"`
        Subcategory string `json:"subCategory" schema:"subcategory"`
        Brand string `json:"brand" schema:"brand"`
        Desc string `json:"desc" schema:"desc"`
        Price float64 `json:"price,string" schema:"price"`
        Ccy string `json:"ccy"`
        Tax float64 `json:"tax,string" schema:"tax"`//percent
        Stock int `json:"stock,string" schema:"stock"`
        Url string `json:"url" schema:"url"`
        Urlimg string `json:"urlImg" schema:"urlImg"`
        Featured bool `json:"featured,string" schema:"featured"`
        Hidden bool `json:"hidden,string" schema:"hidden"`
        DeetId int64
        Deets *GoodDeets `json:"goodDeets" schema:"Deets"`
}
type GoodDeets struct {
        //Id int64 `json:"id"`
        DescDeets string `json:"descDetails" schema:"descDetails"`//path to file with details
        Related []int64 `json:"related" cap:"6" schema:"-"`
        Prices []float64 `json:"prices" cap:"6" schema:"-"`
        Volumes []int `json:"volumes" cap:"6" schema:"-"`
        //PriceVolume map[int]float64 `json:"priceVolume"`
        Parameters []string `json:"parameters" cap:"12" schema:"-"`
        Package []string `json:"package" cap:"12" schema:"-"`
        //Parameters map[string]string `json:"parameters"`
        Features []string `json:"features" cap:"12" schema:"-"`
        Items []string `json:"items" cap:"6" schema:"-"`
        UrlImgs1 string `json:"urlImgs1" schema:"urlImgs"`
        UrlImgs2 string `json:"urlImgs2" schema:"-"`
        UrlImgs3 string `json:"urlImgs3" schema:"-"`
        UrlFile string `json:"urlFile" schema:"urlFile"`
}

type RelatedGoods struct {
        Gd Good `json:"good,string"`
        Rgds []Good `json:"goods,string"`
}

type Cart struct {
        Id int64 `json:"id"`
        //Items []*Good `json:"items"`
        //Units []int `json:"units"`
        Platform string `json:"platform,omitempty"`
        Items map[Good]int `json:"items"`
        TotalUnits int `json:"totalunits"`
        Ccy string `json:"ccy"`
        Subtotal float64 `json:"subtotal"`
        Total float64 `json:"total"`
        RefID string `json:"refid"`
}

type GoodDatabase interface {

        ListGoods(limit int) ([]*Good, error)

        Add(good *Good) (int64, error)

        Get(id int64) (*Good, error)

        GetGoodsByCategory(cat Category) ([]*Good, error)

        GetGoodsBySubcategory(subcat string) ([]*Good, error)

        Update(good *Good) (int64, error)

        DeleteGood(id int64) error

        //Close() error

}
