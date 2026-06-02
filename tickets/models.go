package tickets

import (
	"encoding/json"
	"fmt"
	"time"
)

var statusToPipeline = map[string]string{
	"inbox":              "1",
	"wachten op reactie": "2",
	"terugkoppelen":      "3",
	"complete":           "4",
	"to do":              "5302752480",
	"in progress":        "5302752481",
	"opnieuw geopend":    "5316325601",
}

type UpdateTicketWebhookRequest struct {
	AutoID    string    `json:"auto_id"`
	TriggerID string    `json:"trigger_id"`
	Date      time.Time `json:"date"`
	Payload   struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		Content          string `json:"content"`
		LowerName        string `json:"lower_name"`
		HTMLContent      any    `json:"html_content"`
		TextContent      string `json:"text_content"`
		ContentSize      string `json:"content_size"`
		SprintPoints     any    `json:"sprint_points"`
		Coverimage       any    `json:"coverimage"`
		Priority         any    `json:"priority"`
		PersonalPriority []any  `json:"personal_priority"`
		DraftUUID        any    `json:"draft_uuid"`
		CustomID         any    `json:"custom_id"`
		CustomType       any    `json:"custom_type"`
		StatusID         string `json:"status_id"`
		WorkspaceID      string `json:"workspace_id"`
		Subcategory      string `json:"subcategory"`
		DirectParent     any    `json:"direct_parent"`
		RootParent       any    `json:"root_parent"`
		MergedTo         any    `json:"merged_to"`
		SubtaskSort      any    `json:"subtask_sort"`
		SubtaskSortDir   any    `json:"subtask_sort_dir"`
		Reccurence       struct {
			V1 struct {
				Recurring         bool `json:"recurring"`
				RecurType         any  `json:"recur_type"`
				RecurNext         any  `json:"recur_next"`
				RecurNewStatus    any  `json:"recur_new_status"`
				RecurDueDate      any  `json:"recur_due_date"`
				RecurData         any  `json:"recur_data"`
				RecurRule         any  `json:"recur_rule"`
				RecurTask         any  `json:"recur_task"`
				RecurSkipMissed   bool `json:"recur_skip_missed"`
				RecurOnStatus     any  `json:"recur_on_status"`
				RecurOn           any  `json:"recur_on"`
				RecurCopyOriginal bool `json:"recur_copy_original"`
				RecurTime         bool `json:"recur_time"`
				RecurImmediately  bool `json:"recur_immediately"`
				RecurUntil        any  `json:"recur_until"`
				RecurDst          any  `json:"recur_dst"`
				RecurTzOffset     any  `json:"recur_tz_offset"`
				RecurTz           any  `json:"recur_tz"`
				RecurDaily        bool `json:"recur_daily"`
				RecurIgnoreToday  bool `json:"recur_ignore_today"`
			} `json:"v1"`
			V2 struct {
				SetTime                   bool `json:"set_time"`
				CreateNewTask             bool `json:"create_new_task"`
				Periodically              bool `json:"periodically"`
				SimpleSettings            bool `json:"simple_settings"`
				IgnoreWeekends            bool `json:"ignore_weekends"`
				RecurOnSchedule           bool `json:"recur_on_schedule"`
				RescheduleOnDueDateChange bool `json:"reschedule_on_due_date_change"`
			} `json:"v2"`
		} `json:"reccurence"`
		Privacy struct {
			Private               bool `json:"private"`
			Public                bool `json:"public"`
			PublicToken           any  `json:"public_token"`
			PublicPermissionLevel any  `json:"public_permission_level"`
			PublicFields          any  `json:"public_fields"`
			PublicShareExpiresOn  any  `json:"public_share_expires_on"`
			PublicSharing         bool `json:"public_sharing"`
			SeoOptimized          bool `json:"seo_optimized"`
			MadePublicBy          any  `json:"made_public_by"`
			MadePublicTime        any  `json:"made_public_time"`
		} `json:"privacy"`
		Templating struct {
			Template            bool `json:"template"`
			OriginalSubcat      any  `json:"original_subcat"`
			TemplateName        any  `json:"template_name"`
			TeamID              any  `json:"team_id"`
			TemplateFieldIds    any  `json:"template_field_ids"`
			PermanentTemplateID any  `json:"permanent_template_id"`
			Visibility          any  `json:"visibility"`
		} `json:"templating"`
		States struct {
			IsArchived  bool `json:"is_archived"`
			IsDeleted   bool `json:"is_deleted"`
			IsEncrypted bool `json:"is_encrypted"`
		} `json:"states"`
		TimeMgmt struct {
			IsSummaryTask      bool   `json:"is_summary_task"`
			StartDate          any    `json:"start_date"`
			StartDateTime      bool   `json:"start_date_time"`
			DateClosed         any    `json:"date_closed"`
			DateCreated        string `json:"date_created"`
			DateUpdated        string `json:"date_updated"`
			DueDate            any    `json:"due_date"`
			DueDateTime        bool   `json:"due_date_time"`
			DateDeleted        any    `json:"date_deleted"`
			DateActive         any    `json:"date_active"`
			DateUnstarted      any    `json:"date_unstarted"`
			DateDone           any    `json:"date_done"`
			DateDelegated      any    `json:"date_delegated"`
			TimeEstimate       any    `json:"time_estimate"`
			TimeEstimateString any    `json:"time_estimate_string"`
			TimeSpent          any    `json:"time_spent"`
			SentDueDateNotif   bool   `json:"sent_due_date_notif"`
			Duration           struct {
			} `json:"duration"`
			DurationIsElapsed bool `json:"duration_is_elapsed"`
		} `json:"time_mgmt"`
		Checklists []any `json:"checklists"`
		Ownership  struct {
			CreatedByEmail any    `json:"created_by_email"`
			Owner          int    `json:"owner"`
			Creator        int    `json:"creator"`
			DeletedBy      any    `json:"deleted_by"`
			Source         string `json:"source"`
			Delegator      any    `json:"delegator"`
			FormID         any    `json:"form_id"`
			MergedTo       any    `json:"merged_to"`
		} `json:"ownership"`
		SubtaskIds []any `json:"subtask_ids"`
		Tags       []any `json:"tags"`
		Fields     []struct {
			FieldID       string     `json:"field_id"`
			Value         FlexString `json:"value"`
			ValueOptions  any        `json:"value_options"`
			ValueDeleted  bool       `json:"value_deleted"`
			ValueRichtext any        `json:"value_richtext"`
			Reasoning     any        `json:"reasoning"`
			TaskID        string     `json:"task_id"`
			Type          int        `json:"type"`
		} `json:"fields"`
		Lists []struct {
			ListID string `json:"list_id"`
			Type   string `json:"type"`
		} `json:"lists"`
		Users []struct {
			Userid int    `json:"userid"`
			Type   string `json:"type"`
		} `json:"users"`
		Groups            []any `json:"groups"`
		RelatedTasks      []any `json:"related_tasks"`
		RelatedTasksCount struct {
			DependsOn  int `json:"dependsOn"`
			DependedBy int `json:"dependedBy"`
		} `json:"related_tasks_count"`
		Attachments []struct {
			AttachmentID string `json:"attachment_id"`
			Type         string `json:"type"`
		} `json:"attachments"`
		Docs          []any `json:"docs"`
		VersionVector struct {
			WorkspaceID int    `json:"workspace_id"`
			ObjectType  string `json:"object_type"`
			ObjectID    string `json:"object_id"`
			Vector      []struct {
				MasterID int   `json:"master_id"`
				Version  int64 `json:"version"`
				Deleted  bool  `json:"deleted"`
			} `json:"vector"`
		} `json:"_version_vector"`
	} `json:"payload"`
}
type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexString(s)
		return nil
	}
	// Try number
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexString(n.String())
		return nil
	}
	return fmt.Errorf("cannot unmarshal %s into FlexString", data)
}
