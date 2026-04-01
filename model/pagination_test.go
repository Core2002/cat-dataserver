package model

import (
	"testing"
)

func TestPaginationRequestGetPage(t *testing.T) {
	tests := []struct {
		name     string
		input    PaginationRequest
		expected int
	}{
		{"Zero page defaults to 1", PaginationRequest{Page: 0}, 1},
		{"Negative page defaults to 1", PaginationRequest{Page: -1}, 1},
		{"Valid page", PaginationRequest{Page: 5}, 5},
		{"First page", PaginationRequest{Page: 1}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.GetPage()
			if result != tt.expected {
				t.Errorf("GetPage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPaginationRequestGetPageSize(t *testing.T) {
	tests := []struct {
		name     string
		input    PaginationRequest
		expected int
	}{
		{"Zero size defaults to 10", PaginationRequest{PageSize: 0}, 10},
		{"Negative size defaults to 10", PaginationRequest{PageSize: -5}, 10},
		{"Valid size", PaginationRequest{PageSize: 20}, 20},
		{"Default size 10", PaginationRequest{PageSize: 10}, 10},
		{"Max size 100", PaginationRequest{PageSize: 100}, 100},
		{"Size greater than 100 caps at 100", PaginationRequest{PageSize: 150}, 100},
		{"Size greater than 100 caps at 100 with 200", PaginationRequest{PageSize: 200}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.GetPageSize()
			if result != tt.expected {
				t.Errorf("GetPageSize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPaginationRequestGetOffset(t *testing.T) {
	tests := []struct {
		name     string
		input    PaginationRequest
		expected int
	}{
		{"First page offset", PaginationRequest{Page: 1, PageSize: 10}, 0},
		{"Second page offset", PaginationRequest{Page: 2, PageSize: 10}, 10},
		{"Third page offset", PaginationRequest{Page: 3, PageSize: 20}, 40},
		{"Large page offset", PaginationRequest{Page: 10, PageSize: 50}, 450},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.GetOffset()
			if result != tt.expected {
				t.Errorf("GetOffset() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewPaginationResponse(t *testing.T) {
	tests := []struct {
		name        string
		data        interface{}
		total       int64
		page        int
		pageSize    int
		expectedRes PaginationResponse
	}{
		{
			name:     "Single page",
			data:     []int{1, 2, 3},
			total:    3,
			page:     1,
			pageSize: 10,
			expectedRes: PaginationResponse{
				Data:       []int{1, 2, 3},
				Total:      3,
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
			},
		},
		{
			name:     "Multiple pages",
			data:     []int{1, 2, 3, 4, 5},
			total:    15,
			page:     1,
			pageSize: 5,
			expectedRes: PaginationResponse{
				Data:       []int{1, 2, 3, 4, 5},
				Total:      15,
				Page:       1,
				PageSize:   5,
				TotalPages: 3,
			},
		},
		{
			name:     "Exact division",
			data:     []int{1, 2, 3, 4, 5},
			total:    10,
			page:     2,
			pageSize: 5,
			expectedRes: PaginationResponse{
				Data:       []int{1, 2, 3, 4, 5},
				Total:      10,
				Page:       2,
				PageSize:   5,
				TotalPages: 2,
			},
		},
		{
			name:     "Empty data",
			data:     []int{},
			total:    0,
			page:     1,
			pageSize: 10,
			expectedRes: PaginationResponse{
				Data:       []int{},
				Total:      0,
				Page:       1,
				PageSize:   10,
				TotalPages: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPaginationResponse(tt.data, tt.total, tt.page, tt.pageSize)
			if result.Total != tt.expectedRes.Total {
				t.Errorf("Total = %v, want %v", result.Total, tt.expectedRes.Total)
			}
			if result.Page != tt.expectedRes.Page {
				t.Errorf("Page = %v, want %v", result.Page, tt.expectedRes.Page)
			}
			if result.PageSize != tt.expectedRes.PageSize {
				t.Errorf("PageSize = %v, want %v", result.PageSize, tt.expectedRes.PageSize)
			}
			if result.TotalPages != tt.expectedRes.TotalPages {
				t.Errorf("TotalPages = %v, want %v", result.TotalPages, tt.expectedRes.TotalPages)
			}
		})
	}
}

func TestPaginationResponseTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		pageSize int
		expected int
	}{
		{"Zero total", 0, 10, 0},
		{"One page", 5, 10, 1},
		{"Exact multiple", 20, 10, 2},
		{"Partial last page", 25, 10, 3},
		{"Large partial", 101, 10, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewPaginationResponse(nil, tt.total, 1, tt.pageSize)
			if response.TotalPages != tt.expected {
				t.Errorf("TotalPages = %v, want %v", response.TotalPages, tt.expected)
			}
		})
	}
}

func TestCatEventTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value CatEventType
	}{
		{"CatSick", CatSick},
		{"CatInjure", CatInjure},
		{"CatPreg", CatPreg},
		{"CatBirth", CatBirth},
		{"CatDeath", CatDeath},
		{"CatContractTerminatio", CatContractTerminatio},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Expected non-empty value for %s", tt.name)
			}
		})
	}
}

func TestCatActionTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value CatActionType
	}{
		{"CatActionFeed", CatActionFeed},
		{"CatActionGiveWater", CatActionGiveWater},
		{"CatActionTakeTemperature", CatActionTakeTemperature},
		{"CatActionPlay", CatActionPlay},
		{"CatActionSterilize", CatActionSterilize},
		{"CatActionHealthCheck", CatActionHealthCheck},
		{"CatActionDeworm", CatActionDeworm},
		{"CatActionCleanLitter", CatActionCleanLitter},
		{"CatActionDisinfect", CatActionDisinfect},
		{"CatActionTrimNails", CatActionTrimNails},
		{"CatActionBathing", CatActionBathing},
		{"CatActionVaccinate", CatActionVaccinate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("Expected non-empty value for %s", tt.name)
			}
		})
	}
}

func TestTimeNow(t *testing.T) {
	now := TimeNow()
	if now.IsZero() {
		t.Error("TimeNow() returned zero time")
	}
}
