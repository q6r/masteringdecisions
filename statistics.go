package main

// TODO : Finish implementation
// TODO : Add stats to database

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/montanaflynn/stats"
)

type BallotStatistics struct {
	Criterion_ID     int     `json:"criterion_id"`
	Criterion_Name   string  `json:"criterion_name"`
	Criterion_Weight int     `json:"criterion_weight"`
	Voted            int     `json:"voted"`
	Perc             float64 `json:"perc"`
}

type BallotAnalysis struct {
	Decision_ID                           int                `json:"decision_id"`
	Ballot_ID                             int                `json:"ballot_id"`
	Ballot_Name                           string             `json:"ballot_name"`
	Statistics                            []BallotStatistics `json:"statistics"`
	MinVote                               float64            `json:"min_vote"`
	MaxVote                               float64            `json:"max_vote"`
	AvgVote                               float64            `json:"avg_vote"`
	MeanVote                              float64            `json:"mean_vote"`
	MedianVote                            float64            `json:"median_vote"`
	ModeVote                              float64
	PopulationVarianceVote                float64
	SampleVarianceVote                    float64
	MedianAbsoluteDeviationPopulationVote float64
	StandardDeviationPopulationVote       float64
	StandardDeviationSampleVote           float64

	//...etc
}

// HStats does math analysis/statistics on all ballots beloning to a decision
// it returns an array of BallotAnalysis that the frontend can use to display
// things.
// Note : This will fail if a ballot didn't vote on a criterion "missing ballot"
func HStats(c *gin.Context) {
	did := c.Param("decision_id")

	// Get the decision
	var decision Decision
	err := dbmap.SelectOne(&decision,
		"select * from decision where decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Get ballots for that decision
	var ballots []Ballot
	_, err = dbmap.Select(&ballots,
		"select * from ballot where decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Get criterions for that decisions
	var criterions []Criterion
	_, err = dbmap.Select(&criterions,
		"select * from criterion where decision_id=$1", did)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// find votes for this decision by ballot
	var all_votes [][]Vote
	for _, ballot := range ballots {
		var votes []Vote
		_, err = dbmap.Select(&votes,
			"select * from vote where ballot_id=$1",
			ballot.Ballot_ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		all_votes = append(all_votes, votes)
	}

	// At this points we have
	// 1. The decision
	// 2. Its ballots
	// 3. Its criterions
	// 4. Its votes grouped by (by ballots)

	ballots_analysis := make([]*BallotAnalysis, 0)

	for _, ballot := range ballots {
		ba := &BallotAnalysis{}
		ba.Decision_ID = decision.Decision_ID
		ba.Ballot_ID = ballot.Ballot_ID
		ba.Ballot_Name = ballot.Name

		for _, criterion := range criterions {
			stat := BallotStatistics{}
			stat.Criterion_ID = criterion.Criterion_ID
			stat.Criterion_Name = criterion.Name
			stat.Criterion_Weight = criterion.Weight

			vote, err := FindVotesByKeys(criterion.Criterion_ID, ballot.Ballot_ID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": fmt.Errorf("missing vote for criterion id %d ballot id %d",
						criterion.Criterion_ID, ballot.Ballot_ID).Error(),
				})
				return
			}

			stat.Voted = vote.Weight
			stat.Perc = (float64(stat.Voted) / float64(stat.Criterion_Weight)) * 100
			ba.Statistics = append(ba.Statistics, stat)
		}
		ballots_analysis = append(ballots_analysis, ba)
	}

	// Run the math on the ballots
	for _, ba := range ballots_analysis {

		ba.MaxVote, err = ApplyFunctionOnVotes(ba, stats.Max)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ba.MinVote, err = ApplyFunctionOnVotes(ba, stats.Min)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ba.AvgVote, err = ApplyFunctionOnVotes(ba, stats.Sum)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ba.AvgVote /= float64(len(ba.Statistics))

		ba.MeanVote, err = ApplyFunctionOnVotes(ba, stats.Mean)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ba.MedianVote, err = ApplyFunctionOnVotes(ba, stats.Median)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, ballots_analysis)
}

func ApplyFunctionOnVotes(ba *BallotAnalysis, f func(stats.Float64Data) (float64, error)) (float64, error) {
	var votes []float64
	var err error

	for _, s := range ba.Statistics {
		votes = append(votes, float64(s.Voted))
	}

	result, err := f(votes)
	if err != nil {
		return 0, err
	}

	return result, nil
}
