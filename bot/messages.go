package bot

import (
	"errors"
)

var (
	ErrorNoSenderIdFound = errors.New("no sender ID found in message")
)

const (
	MESSAGE_FORWARD_CONNECTION_01              = "🤝 %s wants to negotiate with you!\nTo join the conversation, send this message to @NegotiMateBot:\n\n/connect %v"
	MESSAGE_FORWARD_CONNECTION_02              = "📬 Forward the message above ☝️ to the person you want to negotiate with"
	MESSAGE_WAITING_FOR_CONNECTION             = "⏳ Waiting for the other person to connect..."
	MESSAGE_YOU_CANT_CONNECT_TO_YOURSELF       = "😅 Oops! You can't negotiate with yourself"
	MESSAGE_NO_SUCH_USER_IS_AWATING            = "❌ This negotiation ID doesn't exist or has expired"
	MESSAGE_START_GUIDE                        = "👋 Welcome to NegotiMate!\n\nHere's how to use the bot:\n🚀 /start - Start a new negotiation\n🤝 /connect <ID> - Join an existing negotiation\n🔄 /reset - Start over\n\n💡 Note: This bot helps you agree on the number only. Remember to discuss other details (currency, bonuses, gross/net, monthly/yearly) separately!"
	MESSAGE_SELECT_YOUR_ROLE_CONNECTED         = "🎉 You're connected! Now, let's get started.\nFirst, select your role"
	MESSAGE_SELECT_YOUR_ROLE_UNEXPECTED        = "🤔 Something unexpected happened...\nLet's try again - please select your role"
	MESSAGE_UNEXPECTED_STATE                   = "😮 Oops! Something went wrong.\nTry starting over with /reset"
	MESSAGE_SELECT_SALARY_LOWER_BOUND_EMPLOYEE = "👔 Your role: %s\nWhat's the minimum salary you'd accept?"
	MESSAGE_SELECT_SALARY_LOWER_BOUND_EMPLOYER = "👔 Your role: %s\nWhat's the minimum salary you can offer?"
	MESSAGE_SELECT_SALARY_UPPER_BOUND_EMPLOYEE = "💰 Now, what's the maximum salary you're willing to consider?"
	MESSAGE_SELECT_SALARY_UPPER_BOUND_EMPLOYER = "💰 Now, what's the maximum salary you're can offer?"
	MESSAGE_WAITING_FOR_RESULT                 = "🎲 Crunching the numbers..."
	MESSAGE_USE_VALID_POSITIVE_NUMBER          = "📊 Please enter a valid positive number"
	MESSAGE_USE_VALID_UPPER_BOUND              = "⚠️ The maximum should be no more than %d times your minimum"
	MESSAGE_RESULT_SUCCESS                     = "🎊 Great news! You both can agree on: %d"
	MESSAGE_RESULT_ERROR                       = "😕 Unfortunately, your salary ranges don't overlap. Want to try again?"
)
