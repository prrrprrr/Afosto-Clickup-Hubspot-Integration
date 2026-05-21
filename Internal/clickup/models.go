package clickup

const (
	createTaskEndPoint         = "https://api.clickup.com/api/v2/list/901521484881/task"  // create task endpoint of destination list
	getCustomFieldEndPoint     = "https://api.clickup.com/api/v2/list/901521484881/field" // custom field endpoint of destination list
	clickupCommentEndpoint     = "https://api.clickup.com/api/v2/task/%s/comment"         // comment endpoint, needs ClickUpTaskID
	attachmentEndpoint         = "https://api.clickup.com/api/v2/task/%s/attachment"
	EmailContactCustomfield    = "e67b5243-eeed-49b7-ba72-76d48c040515"
	HubSpotTicketIdCustomField = "e210026c-aefd-432d-9778-5673b7e1543a"
	ThreadedCommentEndpoint    = "https://api.clickup.com/api/v2/comment/%s/reply"
	ClickupTaskEndpoint        = "https://api.clickup.com/api/v2/task/%s"
	moveTaskEndpoint           = "https://api.clickup.com/api/v3/workspaces/1373770/tasks/%s/home_list/%s" // url for moving a list, first %s is for the task id and the 2nd is for the list id
	sprintFolderID             = "90153990938"
	getListsEndpoint           = "https://api.clickup.com/api/v2/folder/%s/list"
	CompleteStatus             = "complete"
	WachtenOpReactieStatus     = "wachten op reactie"
	InProgressStatus           = "in progress"
	TerugkoppelenStatus        = "terugkoppelen"
	InboxStatus                = "inbox"
	ToDoStatus                 = "to do"
	OpnieuwGeopendStatus       = "opnieuw geopend"
	ReactieOntvangenStatus     = "reactie ontvangen"
)
