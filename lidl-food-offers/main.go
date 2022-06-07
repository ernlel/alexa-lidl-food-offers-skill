package main

import (
	"alexa-lidl-offers-skill/alexa"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler(request alexa.Request) (alexa.Response, error) {
	return IntentSwitch(request)
}

// switch intent handlers
func IntentSwitch(request alexa.Request) (alexa.Response, error) {
	var (
		response alexa.Response
		err      error
	)

	switch request.Body.Intent.Name {
	case "LidlDealsIntent":
		response, err = HandleLidlFoodOffersIntent(request)
	case alexa.HelpIntent:
		response = HandleHelpIntent(request)
	case "AboutIntent":
		response = HandleAboutIntent(request)
	default:
		response, err = HandleLidlFoodOffersIntent(request)
	}
	if err != nil {
		return response, err
	}
	return response, nil
}

// About intent
func HandleAboutIntent(request alexa.Request) alexa.Response {
	return alexa.NewSimpleResponse("About", "Lidl offers was created by Ernestas Leliuga as an unofficial Lidl (UK) food offers Alexa skill.")
}

// Help intent
func HandleHelpIntent(request alexa.Request) alexa.Response {
	var builder alexa.SSMLBuilder
	builder.Say("Here are some of the things you can ask:")
	builder.Pause("1000")
	builder.Say("Alexa, Tell lidl food offers")
	builder.Pause("1000")
	builder.Say("Alexa, lidl food offers")
	return alexa.NewSSMLResponse("Lidl food offers help", builder.Build())
}

// Lidl food offers intent
func HandleLidlFoodOffersIntent(request alexa.Request) (alexa.Response, error) {
	offers, err := getOffers()
	if err != nil {
		return alexa.Response{}, err
	}
	var builder alexa.SSMLBuilder
	builder.Say("Here are Lidl food offers:")
	builder.Pause("1000")
	for _, offer := range offers {
		builder.Say(offer.name + " at " + offer.discount + " percent discount. Now just £" + offer.offerPrice)
		builder.Pause("1000")
	}
	return alexa.NewSSMLResponse("Lidl offers", builder.Build()), nil
}
