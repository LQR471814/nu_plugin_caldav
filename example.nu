let homeset = caldav query homeset
let calendars = caldav query calendars $homeset
let cal = $calendars | filter { |it| $it.name == "Schedule" } | $in.0.path
caldav query events $cal | select uid name categories | filter { |it| ($it.categories | length) == 0 }
