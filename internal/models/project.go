//
// @project GeniusRabbit rotator 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package models

import "geniusrabbit.dev/sspserver/internal/billing"

// Project model
type Project struct {
	ID           uint64        // Authoincrement key
	UserID       uint64        //
	Balance      billing.Money //
	MaxDaily     billing.Money //
	Spent        billing.Money //
	RevenueShare float64       // From 0 to 100 percents
}

// // ProjectFromModel convert database model to specified model
// func ProjectFromModel(p models.Project) Project {
// 	return Project{
// 		ID: p.ID,
// 	}
// }

// RevenueShareFactor amount %
func (p *Project) RevenueShareFactor() float64 {
	return p.RevenueShare / 100.0
}

// ComissionShareFactor which system get from publisher 0..1
func (p *Project) ComissionShareFactor() float64 {
	return (100.0 - p.RevenueShare) / 100.0
}
