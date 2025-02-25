package service

import (
	pageBusiness "crave/miner/cmd/api/domain/business/page"
	"crave/miner/cmd/api/infrastructure/externalApi"
	"crave/miner/cmd/api/infrastructure/repository"
	"crave/miner/cmd/model"

	craveModel "crave/shared/model"
)

type Service struct {
	pageStrat *pageBusiness.PageStrategy
	repo      repository.IRepository
	hubClient externalApi.IHubClient
}

func NewService(pageStrat *pageBusiness.PageStrategy, repo repository.IRepository, hubClient externalApi.IHubClient) *Service {
	return &Service{pageStrat: pageStrat, repo: repo, hubClient: hubClient}
}

func (s *Service) Parse(step craveModel.Step, page craveModel.Page, name string, filter craveModel.Filter) error {
	pageBiz := s.getPageBusiness(page)
	targets, err := s.getNextTargets(pageBiz, step, name)
	if err != nil {
		return err
	}
	s.repo.Save(name, page, targets)
	s.hubClient.ParseResult(name, s.getNames(targets), step)
	return nil
}

func (s *Service) getNames(targets []model.ParsedTarget) []string {
	targetNames := make([]string, len(targets))
	for i, target := range targets {
		targetNames[i] = target.Name
	}
	return targetNames
}

func (s *Service) getPageBusiness(page craveModel.Page) pageBusiness.IBusiness {
	return s.pageStrat.GetPageBusiness(page)
}

func (s *Service) getNextTargets(page pageBusiness.IBusiness, step craveModel.Step, name string) ([]model.ParsedTarget, error) {
	return page.ParseNextTargets(step, name)
}
