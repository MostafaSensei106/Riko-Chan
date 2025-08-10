package bot

import "github.com/MostafaSensei106/Riko-Chan/internal/models"

// getText returns localized text based on user language
func (b *Bot) getText(key string, language models.UserLanguage) string {
	texts := map[models.UserLanguage]map[string]string{
		models.LanguageEnglish: {
			"welcome":               "🌟 Welcome to Future Message Bot! 🌟\n\nI help you schedule messages to be sent in the future.",
			"help_text":             "Use /new to create a message, /list to view pending messages, and /settings to configure your preferences.",
			"new_message":           "📝 New Message",
			"my_messages":           "📋 My Messages",
			"settings":              "⚙️ Settings",
			"help":                  "❓ Help",
			"unknown_command":       "Unknown command. Type /help to see available commands.",
			"new_message_help":      "Usage: /new <message> at <time>\nExample: /new Hello world at 2024-01-01 15:30",
			"invalid_format":        "Invalid format. Use: <message> at <time>",
			"invalid_time_format":   "Invalid time format. Examples: 'tomorrow 9:00', 'after 2 hours', '2024-01-01 15:30'",
			"error_occurred":        "An error occurred. Please try again.",
			"message_scheduled":     "✅ Message scheduled for %s\n🆔 ID: %s",
			"add_notification":      "🔔 Add Notification",
			"make_recurring":        "🔄 Make Recurring",
			"send_to_other":         "👤 Send to Other",
			"no_pending_messages":   "You have no pending messages.",
			"your_pending_messages": "📋 Your Pending Messages:",
			"message":               "Message",
			"cancel_help":           "Usage: /cancel <message_id>",
			"delete_help":           "Usage: /delete <message_id>",
			"message_not_found":     "Message not found.",
			"message_cancelled":     "Message cancelled successfully.",
			"message_deleted":       "Message deleted successfully.",
			"change_language":       "Change Language",
			"change_timezone":       "Change Timezone",
			"integrations":          "Integrations",
			"current_settings":      "🛠 Current Settings:\n🌍 Language: %s\n🕒 Timezone: %s",
			"detailed_help":         "🤖 Future Message Bot Help\n\n📝 Commands:\n/new <message> at <time> - Schedule a message\n/list - View pending messages\n/cancel <id> - Cancel a message\n/delete <id> - Delete a message\n/settings - Configure settings\n\n⏰ Time formats:\n- 'after 2 hours'\n- 'tomorrow 9:00'\n- '2024-01-01 15:30'\n- 'next Friday 14:00'",
			"new_message_prompt":    "Please send your message in the format:\n<message> at <time>",
			"unclear_message":       "I didn't understand. Use /help to see how to use me.",
		},
		models.LanguageArabic: {
			"welcome":               "🌟 أهلاً بك في بوت الرسائل المستقبلية! 🌟\n\nأساعدك في جدولة الرسائل لإرسالها في المستقبل.",
			"help_text":             "استخدم /new لإنشاء رسالة، /list لعرض الرسائل المعلقة، و /settings لتكوين تفضيلاتك.",
			"new_message":           "📝 رسالة جديدة",
			"my_messages":           "📋 رسائلي",
			"settings":              "⚙️ الإعدادات",
			"help":                  "❓ المساعدة",
			"unknown_command":       "أمر غير معروف. اكتب /help لرؤية الأوامر المتاحة.",
			"new_message_help":      "الاستخدام: /new <الرسالة> at <الوقت>\nمثال: /new مرحبا بالعالم at 2024-01-01 15:30",
			"invalid_format":        "تنسيق غير صحيح. استخدم: <الرسالة> at <الوقت>",
			"invalid_time_format":   "تنسيق وقت غير صحيح. أمثلة: 'غداً 9:00'، 'بعد ساعتين'، '2024-01-01 15:30'",
			"error_occurred":        "حدث خطأ. يرجى المحاولة مرة أخرى.",
			"message_scheduled":     "✅ تم جدولة الرسالة لـ %s\n🆔 المعرف: %s",
			"add_notification":      "🔔 إضافة تنبيه",
			"make_recurring":        "🔄 جعلها متكررة",
			"send_to_other":         "👤 إرسال لشخص آخر",
			"no_pending_messages":   "ليس لديك رسائل معلقة.",
			"your_pending_messages": "📋 رسائلك المعلقة:",
			"message":               "رسالة",
			"cancel_help":           "الاستخدام: /cancel <معرف_الرسالة>",
			"delete_help":           "الاستخدام: /delete <معرف_الرسالة>",
			"message_not_found":     "الرسالة غير موجودة.",
			"message_cancelled":     "تم إلغاء الرسالة بنجاح.",
			"message_deleted":       "تم حذف الرسالة بنجاح.",
			"change_language":       "تغيير اللغة",
			"change_timezone":       "تغيير المنطقة الزمنية",
			"integrations":          "التكاملات",
			"current_settings":      "🛠 الإعدادات الحالية:\n🌍 اللغة: %s\n🕒 المنطقة الزمنية: %s",
			"detailed_help":         "🤖 مساعدة بوت الرسائل المستقبلية\n\n📝 الأوامر:\n/new <رسالة> at <وقت> - جدولة رسالة\n/list - عرض الرسائل المعلقة\n/cancel <معرف> - إلغاء رسالة\n/delete <معرف> - حذف رسالة\n/settings - تكوين الإعدادات\n\n⏰ تنسيقات الوقت:\n- 'بعد ساعتين'\n- 'غداً 9:00'\n- '2024-01-01 15:30'\n- 'الجمعة القادمة 14:00'",
			"new_message_prompt":    "يرجى إرسال رسالتك بالتنسيق:\n<الرسالة> at <الوقت>",
			"unclear_message":       "لم أفهم. استخدم /help لمعرفة كيفية استخدامي.",
		},
	}

	if langTexts, exists := texts[language]; exists {
		if text, exists := langTexts[key]; exists {
			return text
		}
	}

	// Fallback to English
	if langTexts, exists := texts[models.LanguageEnglish]; exists {
		if text, exists := langTexts[key]; exists {
			return text
		}
	}

	return key // Return key if no translation found
}
