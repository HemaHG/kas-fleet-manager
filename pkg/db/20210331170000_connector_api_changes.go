package db

// Migrations should NEVER use types from other packages. Types can change
// and then migrations run on a _new_ database will fail or behave unexpectedly.
// Instead of importing types, always re-create the type in the migration, as
// is done here, even though the same type is defined in pkg/api

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func connectorApiChanges() *gormigrate.Migration {

	actions := []func(tx *gorm.DB, do bool) error{

		func(tx *gorm.DB, do bool) error {
			if do {
				return tx.Exec("ALTER TABLE connector_clusters RENAME COLUMN status TO status_phase").Error
			} else {
				return tx.Exec("ALTER TABLE connector_clusters RENAME COLUMN status_phase TO status").Error
			}
		},

		func(tx *gorm.DB, do bool) error {
			type ConnectorClusterStatus struct {
				Version    string `json:"version,omitempty"`
				Conditions string `json:"conditions,omitempty"`
				Operators  string `json:"operators,omitempty"`
			}
			type ConnectorCluster struct {
				Status ConnectorClusterStatus `json:"status" gorm:"embedded;embedded_prefix:status_"`
			}
			if do {
				return tx.AutoMigrate(&ConnectorCluster{}).Error
			} else {
				if err := tx.Table("connector_clusters").DropColumn("status_version").Error; err != nil {
					return err
				}
				if err := tx.Table("connector_clusters").DropColumn("status_conditions").Error; err != nil {
					return err
				}
				if err := tx.Table("connector_clusters").DropColumn("status_operators").Error; err != nil {
					return err
				}
				return nil
			}
		},

		func(tx *gorm.DB, do bool) error {
			if do {
				return tx.Exec("ALTER TABLE connectors ALTER COLUMN connector_spec TYPE json USING to_json(connector_spec)").Error
			} else {
				return tx.Exec("ALTER TABLE connectors ALTER COLUMN connector_spec TYPE text").Error
			}
		},

		func(tx *gorm.DB, do bool) error {
			if do {
				return tx.Exec("ALTER TABLE connectors RENAME COLUMN addon_group TO addon_cluster_id").Error
			} else {
				return tx.Exec("ALTER TABLE connectors RENAME COLUMN addon_cluster_id TO addon_group").Error
			}
		},

		func(tx *gorm.DB, do bool) error {
			type ConnectorStatus struct {
				Model
				ClusterID string `json:"cluster_id"`
				Phase     string
			}
			type Connector struct {
				Model
				Status ConnectorStatus
			}
			if do {
				if err := tx.AutoMigrate(&ConnectorStatus{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(&Connector{}).Error; err != nil {
					return err
				}
			} else {
				if err := tx.DropTable(&ConnectorStatus{}).Error; err != nil {
					return err
				}
			}
			return nil
		},

		func(tx *gorm.DB, do bool) error {
			type Connector struct {
				Model
				ClusterID string
				Status    string
			}
			if do {
				if err := tx.Table("connectors").DropColumn("status").Error; err != nil {
					return err
				}
				if err := tx.Table("connectors").DropColumn("cluster_id").Error; err != nil {
					return err
				}
			} else {
				if err := tx.AutoMigrate(&Connector{}).Error; err != nil {
					return err
				}
			}
			return nil
		},

		func(tx *gorm.DB, do bool) error {
			type Connector struct {
				DesiredState string `json:"desired_state"`
			}
			if do {
				return tx.AutoMigrate(&Connector{}).Error
			} else {
				return tx.Table("connectors").DropColumn("desired_state").Error
			}
		},

		func(tx *gorm.DB, do bool) error {
			type ConnectorClusters struct {
				AddonGroup string
			}
			if do {
				return tx.Table("connector_clusters").DropColumn("addon_group").Error
			} else {
				return tx.AutoMigrate(&ConnectorClusters{}).Error
			}
		},

		func(tx *gorm.DB, do bool) error {

			type ConnectorDeploymentStatus struct {
				Model
				Phase        string
				SpecChecksum string
				Conditions   string `gorm:"type:jsonb;index:"`
			}

			type ConnectorDeployment struct {
				Model
				Version              int64                     `gorm:"type:bigserial;index:"`
				ConnectorID          string                    `json:"connector_id"`
				ConnectorVersion     int64                     `json:"connector_resource_version"`
				Status               ConnectorDeploymentStatus `gorm:"foreignkey:ID"`
				ConnectorTypeService string
				ClusterID            string `gorm:"index:"`
				SpecChecksum         string `json:"spec_checksum,omitempty"`
			}

			if do {
				if err := tx.AutoMigrate(&ConnectorDeploymentStatus{}).Error; err != nil {
					return err
				}
				if err := tx.AutoMigrate(&ConnectorDeployment{}).Error; err != nil {
					return err
				}
			} else {
				if err := tx.DropTable(&ConnectorDeployment{}).Error; err != nil {
					return err
				}
				if err := tx.DropTable(&ConnectorDeploymentStatus{}).Error; err != nil {
					return err
				}
			}
			return nil
		},

		func(tx *gorm.DB, do bool) error {
			if do {
				return tx.Exec(`
                CREATE FUNCTION connector_deployments_version_trigger() RETURNS TRIGGER LANGUAGE plpgsql AS '
					BEGIN
					NEW.version := nextval(''connector_deployments_version_seq'');
					RETURN NEW;
					END;
				'
			`).Error
			} else {
				return tx.Exec(`DROP FUNCTION connector_deployments_version_trigger`).Error
			}
		},

		func(tx *gorm.DB, do bool) error {
			if do {
				return tx.Exec(`
				CREATE TRIGGER connector_deployments_version_trigger BEFORE INSERT OR UPDATE ON connector_deployments
				FOR EACH ROW EXECUTE PROCEDURE connector_deployments_version_trigger();
			`).Error
			} else {
				return tx.Exec(`DROP TRIGGER connector_deployments_version_trigger ON connector_deployments`).Error
			}
		},
	}

	return &gormigrate.Migration{
		ID: "202103171200",
		Migrate: func(tx *gorm.DB) error {
			for _, action := range actions {
				err := action(tx, true)
				if err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for i := len(actions) - 1; i >= 0; i-- {
				err := actions[i](tx, false)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}