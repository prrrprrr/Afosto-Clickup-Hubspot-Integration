package comments

import "time"

var EmailToActorId = map[string]string{
	"stef@afosto.net":    "88630070",
	"gijs@afosto.com":    "60473959",
	"ton.koop@gmail.com": "23951706",
	"sander@afosto.com":  "12941389",
	"peter@afosto.com":   "67003755",
	"thom@afosto.com":    "30821425",
	"jordin@afosto.com":  "12736422",
	"gertjan@afosto.com": "23951706",
	"ronald@afosto.com":  "23951706",
	"sjoerd@afosto.com":  "23951706",
	"robin@afosto.com":   "23951706",
}

type TaskComments struct {
	Comments []struct {
		ID      string `json:"id,omitempty"`
		Comment []struct {
			Type  string `json:"type,omitempty"`
			Text  string `json:"text,omitempty"`
			Image struct {
				ID              string `json:"id,omitempty"`
				Name            string `json:"name,omitempty"`
				Title           string `json:"title,omitempty"`
				Type            string `json:"type,omitempty"`
				Extension       string `json:"extension,omitempty"`
				ThumbnailLarge  string `json:"thumbnail_large,omitempty"`
				ThumbnailMedium string `json:"thumbnail_medium,omitempty"`
				ThumbnailSmall  string `json:"thumbnail_small,omitempty"`
				URL             string `json:"url,omitempty"`
				Uploaded        bool   `json:"uploaded,omitempty"`
			} `json:"image,omitempty"`
			Attributes struct {
				Width          string `json:"width,omitempty"`
				DataID         string `json:"data-id,omitempty"`
				DataAttachment struct {
					ID              string `json:"id,omitempty"`
					Version         string `json:"version,omitempty"`
					Date            int64  `json:"date,omitempty"`
					Name            string `json:"name,omitempty"`
					Title           string `json:"title,omitempty"`
					Extension       string `json:"extension,omitempty"`
					Source          int    `json:"source,omitempty"`
					ThumbnailSmall  string `json:"thumbnail_small,omitempty"`
					ThumbnailMedium string `json:"thumbnail_medium,omitempty"`
					ThumbnailLarge  string `json:"thumbnail_large,omitempty"`
					Width           int    `json:"width,omitempty"`
					Height          int    `json:"height,omitempty"`
					URL             string `json:"url,omitempty"`
					URLWQuery       string `json:"url_w_query,omitempty"`
					URLWHost        string `json:"url_w_host,omitempty"`
					BlockID         string `json:"block-id,omitempty"`
					List            struct {
						List string `json:"list,omitempty"`
					} `json:"list,omitempty"`
				} `json:"data-attachment,omitempty"`
				DataNaturalWidth  string `json:"data-natural-width,omitempty"`
				DataNaturalHeight string `json:"data-natural-height,omitempty"`
			} `json:"attributes,omitempty"`
		} `json:"comment,omitempty"`
		CommentText string `json:"comment_text,omitempty"`
		User        struct {
			ID             int    `json:"id,omitempty"`
			Username       string `json:"username,omitempty"`
			Email          string `json:"email,omitempty"`
			Color          string `json:"color,omitempty"`
			Initials       string `json:"initials,omitempty"`
			ProfilePicture string `json:"profilePicture,omitempty"`
		} `json:"user,omitempty"`
		Assignee      any    `json:"assignee,omitempty"`
		GroupAssignee any    `json:"group_assignee,omitempty"`
		Reactions     []any  `json:"reactions,omitempty"`
		Date          string `json:"date,omitempty"`
		ReplyCount    int    `json:"reply_count,omitempty"`
	} `json:"comments,omitempty"`
}

type RecieveEmailWebhook struct {
	CallbackID string `json:"callbackId"`
	Origin     struct {
		PortalID                       int `json:"portalId"`
		UserID                         any `json:"userId"`
		ActionDefinitionID             int `json:"actionDefinitionId"`
		ActionDefinitionVersion        int `json:"actionDefinitionVersion"`
		ActionExecutionIndexIdentifier struct {
			EnrollmentID         int64 `json:"enrollmentId"`
			ActionExecutionIndex int   `json:"actionExecutionIndex"`
		} `json:"actionExecutionIndexIdentifier"`
		ExtensionDefinitionID        int `json:"extensionDefinitionId"`
		ExtensionDefinitionVersionID int `json:"extensionDefinitionVersionId"`
	} `json:"origin"`
	Context struct {
		WorkflowID                     int64 `json:"workflowId"`
		ActionID                       int   `json:"actionId"`
		ActionExecutionIndexIdentifier struct {
			EnrollmentID         int64 `json:"enrollmentId"`
			ActionExecutionIndex int   `json:"actionExecutionIndex"`
		} `json:"actionExecutionIndexIdentifier"`
		Source string `json:"source"`
	} `json:"context"`
	Object struct {
		ObjectID   int64  `json:"objectId"`
		ObjectType string `json:"objectType"`
	} `json:"object"`
	Fields struct {
		ClickupTaskID string `json:"clickup_task_id"`
		EmailID       string `json:"email_id"`
		TicketStatus  string `json:"ticket_status"`
		TicketId      string `json:"ticket_id"`
	} `json:"fields"`
	InputFields struct {
		ClickupTaskID string `json:"clickup_task_id"`
		EmailID       string `json:"email_id"`
	} `json:"inputFields"`
	TypedInputs struct {
		ClickupTaskID struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"clickup_task_id"`
		EmailID struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"email_id"`
	} `json:"typedInputs"`
}
type SendEmailWebhookRequest struct {
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
			FieldID       string `json:"field_id"`
			Value         any    `json:"value"`
			ValueOptions  any    `json:"value_options"`
			ValueDeleted  bool   `json:"value_deleted"`
			ValueRichtext any    `json:"value_richtext"`
			Reasoning     any    `json:"reasoning"`
			TaskID        string `json:"task_id"`
			Type          int    `json:"type"`
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
