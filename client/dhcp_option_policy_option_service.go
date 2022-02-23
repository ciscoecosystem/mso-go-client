package client

import (
	"fmt"

	"github.com/ciscoecosystem/mso-go-client/models"
)

func (client *Client) CreateDHCPOptionPolicyOption(obj *models.DHCPOptionPolicyOption) error {
	optionPolicyID, err := client.GetDHCPOptionPolicyID(obj.Name)
	if err != nil {
		return err
	}
	optionPolicyCont, err := client.ReadDHCPOptionPolicy(optionPolicyID)
	if err != nil {
		return err
	}
	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(optionPolicyCont)
	if err != nil {
		return err
	}

	option := models.DHCPOption{
		Data: obj.Data,
		ID:   obj.ID,
		Name: obj.Name,
	}

	DHCPOptionPolicy.DHCPOption = append(DHCPOptionPolicy.DHCPOption, option)
	_, err = client.UpdateDHCPOptionPolicy(optionPolicyID, DHCPOptionPolicy)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) ReadDHCPOptionPolicyOption(obj *models.DHCPOptionPolicyOption) (*models.DHCPOptionPolicyOption, error) {
	optionPolicyID, err := client.GetDHCPOptionPolicyID(obj.Name)
	if err != nil {
		return nil, err
	}
	optionPolicyCont, err := client.ReadDHCPOptionPolicy(optionPolicyID)
	if err != nil {
		return nil, err
	}
	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(optionPolicyCont)
	if err != nil {
		return nil, err
	}

	flag := false
	for _, option := range DHCPOptionPolicy.DHCPOption {
		if option.Name == obj.Name && option.ID == obj.ID {
			flag = true
			break
		}
	}
	if flag {
		return obj, nil
	}
	return nil, fmt.Errorf("No DHCP Option Policy found")
}

func (client *Client) UpdateDHCPOptionPolicyOption(new *models.DHCPOptionPolicyOption, old *models.DHCPOptionPolicyOption) error {
	optionPolicyID, err := client.GetDHCPOptionPolicyID(old.Name)
	if err != nil {
		return err
	}
	optionPolicyCont, err := client.ReadDHCPOptionPolicy(optionPolicyID)
	if err != nil {
		return err
	}
	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(optionPolicyCont)
	if err != nil {
		return err
	}

	NewOptions := make([]models.DHCPOption, 0, 1)
	NewOption := models.DHCPOption{
		Data: new.Data,
		ID:   new.ID,
		Name: new.Name,
	}

	for _, option := range DHCPOptionPolicy.DHCPOption {
		if option.Name != old.Name && option.ID != old.ID {
			NewOptions = append(NewOptions, option)
		} else {
			NewOptions = append(NewOptions, NewOption)
		}
	}

	DHCPOptionPolicy.DHCPOption = NewOptions
	_, err = client.UpdateDHCPOptionPolicy(optionPolicyID, DHCPOptionPolicy)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) DeleteDHCPOptionPolicyOption(obj *models.DHCPOptionPolicyOption) error {
	optionPolicyID, err := client.GetDHCPOptionPolicyID(obj.Name)
	if err != nil {
		return err
	}
	optionPolicyCont, err := client.ReadDHCPOptionPolicy(optionPolicyID)
	if err != nil {
		return err
	}
	DHCPOptionPolicy, err := models.DHCPOptionPolicyFromContainer(optionPolicyCont)
	if err != nil {
		return err
	}
	NewOptions := make([]models.DHCPOption, 0, 1)
	for _, option := range DHCPOptionPolicy.DHCPOption {
		if option.Name != obj.Name && option.ID != obj.ID {
			NewOptions = append(NewOptions, option)
		}
	}
	DHCPOptionPolicy.DHCPOption = NewOptions
	_, err = client.UpdateDHCPOptionPolicy(optionPolicyID, DHCPOptionPolicy)
	if err != nil {
		return err
	}
	return nil
}
