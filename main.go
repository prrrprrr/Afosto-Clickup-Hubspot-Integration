package main

import (
	"Afosto-Clickup-Hubspot-Integration/internal/clickup"
	"Afosto-Clickup-Hubspot-Integration/internal/hubspot"
	"Afosto-Clickup-Hubspot-Integration/attachments"
	"Afosto-Clickup-Hubspot-Integration/comments"
	"Afosto-Clickup-Hubspot-Integration/handlers"
	"Afosto-Clickup-Hubspot-Integration/tasks"
	"Afosto-Clickup-Hubspot-Integration/tickets"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {

	client := http.Client{}

	var hubspotClient hubspot.HubSpotClient

	{
		settings := hubspot.Settings{
			Client: &client,
			APIKey: os.Getenv("HUBSPOT_ACCESS_TOKEN"),
		}

		hubspotClient = hubspot.NewClient(settings)

	}

	var clickupClient clickup.ClickUpClient

	{
		clickUpSettings := clickup.Settings{
			Client:        &client,
			ClickupApiKey: os.Getenv("CLICKUP_ACCESS_TOKEN"),
			HubSpotAPiKey: os.Getenv("HUBSPOT_ACCESS_TOKEN"),
		}

		clickupClient = clickup.NewClient(clickUpSettings)
	}

	var attachmentService attachments.AttachmentService

	{
		attachmentService = attachments.NewService(clickupClient, hubspotClient)
	}

	var ticketService tickets.TicketService

	{
		ticketService = tickets.NewService(hubspotClient, clickupClient)
	}

	var commentService comments.CommentService

	{
		commentService = comments.NewService(comments.Settings{
			ClickupClient:     clickupClient,
			HubspotClient:     hubspotClient,
			AttachmentService: attachmentService,
			ActorID:           os.Getenv("ACTOR_ID"),
		})
	}

	var taskService tasks.TaskService

	{
		taskService = tasks.NewService(tasks.Settings{
			ClickupClient:              clickupClient,
			HubspotClient:              hubspotClient,
			AttachmentService:          attachmentService,
			HubSpotTicketCustomFieldID: os.Getenv("CLICKUP_CUSTOM_FIELD_ID_HUBSPOT_TICKET_ID"),
			CustomerCustomFieldID:      os.Getenv("CLICKUP_CUSTOM_FIELD_ID_CUSTOMER"),
			FallbackCustomerID:         os.Getenv("FALLBACK_CLICKUP_CUSTOMER_ID"),
			HubSpotPortalID:            os.Getenv("HUBSPOT_PORTAL_ID"),
			EmailContactID:             os.Getenv("CLICKUP_CUSTOM_FIELD_ID_EMAIL_CONTACT"),
		})
	}

	var taskHandler = handlers.TaskHandler{Service: taskService}
	var commentHandler = handlers.CommentHandler{Service: commentService}
	var ticketHandler = handlers.TicketHandler{Service: ticketService}

	//setup router
	mux := mux.NewRouter()
	webhook := mux.PathPrefix("/webhook").Subrouter()
	hubspot := webhook.PathPrefix("/hubspot").Subrouter()

	//All calls coming from hubspot are sorted under the subRouter hubspot
	hubspot.HandleFunc("/create_task", handlers.CreateTask(taskHandler)).Methods(http.MethodPost)
	hubspot.HandleFunc("/update_task", handlers.UpdateTask(taskHandler)).Methods(http.MethodPost)
	hubspot.HandleFunc("/receive_email", handlers.ReceiveEmail(commentHandler)).Methods(http.MethodPost)

	//All calls coming from clickup are sorted under the subRouter clickup
	clickup := webhook.PathPrefix("/clickup").Subrouter()
	clickup.HandleFunc("/send_email", handlers.SendEmail(commentHandler)).Methods(http.MethodPost)
	clickup.HandleFunc("/update_ticket", handlers.UpdateTicket(ticketHandler)).Methods(http.MethodPost)

	//run app
	listeningPort := os.Getenv("LISTENING_PORT")
	log.Printf("running on port :%s", listeningPort)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
