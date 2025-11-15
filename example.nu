let principal = caldav query principal
let homeset = caldav query homeset $principal
let calendars = caldav query calendars $homeset
print $calendars

def localtime [str: string] {
	$str | into datetime -z "l"
}

let calpath: string = $calendars
	| filter { |it| $it.name == "Schedule" }
	| first
	| $in.path

caldav query events $calpath | caldav timeline --start (localtime "2025-11-9") --end (localtime "2025-11-16")

[{
	main: {
		summary: "test! 2",
		location: "location of event",
		description: "a human-friendly\ndescription",
		categories: ["category_A", "category_B"],
		class: "PUBLIC", # one of: PUBLIC, PRIVATE, CONFIDENTIAL
		geo: {
			latitude: 0.2,
			longitude: 0.3,
		},
		priority: 1, # 1-9 (1 is highest, 9 is lowest)
		status: "TENTATIVE", # one of: TENTATIVE, CONFIRMED, CANCELLED
		transparency: "", # one of: OPAQUE, TRANSPARENT
		url: "https://github.com",
		comment: "comments that may be relevant",
		attach: "https://github.com/LQR471814",
		contact: "this is contact information associated w/ the event",
		organizer: "https://github.com/LQR41814",
		created: {
			stamp: (date now)
		},
		start: {
			stamp: (localtime "2025-11-14 21:45"),
		},
		end: {
			stamp: (localtime "2025-11-14 22:00"),
		},
	},
}] | caldav save events $calpath

