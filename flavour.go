package main

// This file contains mostly static strings such as names and flavour text
var (
	// functionally constant

	// FirstNamesMale is a list of male names of citizens to be picked from at random.
	FirstNamesMale = []string{
		"Jacob",
		"Thomas",
		"Seth",
		"Isaac",
		"Aidan",
	}
	// FirstNamesFemale
	FirstNamesFemale = []string{
		"Jamie",
		"Rebecca",
		"Donna",
		"Daisy",
		"Lydia",
	}

	// Maybe children should be Name the Second (depending on gender)?
	// That'd be cool and would mean I don't have to have as much variety, wew
	// Considering gender should be random (50/50), we could end up in a situation
	//  where there can be no more births. Could this be mitigated late game by,
	//  I don't know... IVF/gender selection research? Initial settlement citizens
	// 	should be hardcoded with a 50/50 distribution. Now that I think about it,
	// 	without putting the effort in, births will be implicitly incestuous...
	//  TH reckons I should make it a game mechanic, to affect genetics property?
)
