package interfaces

// Signal defines interface for a signal - which a collector + diagnoser will each change to return, from a new separate
// method chain which hooks into the existing "diagnose" / "collect" methods.
type Signal interface {

	//1. getter for data struct representing results.
	//1.1 Mapping of diagnoser/collector name + params data bundle to the output.
	//1.2 Some kind of "priority" indicator on the results, to allow saying "we found something at x level of "importance"
	//    to allow pre-empting / cancelling other signals which we know we don't need.
	//    Conceptually this could be like a "conclusion" that would bubble up to a certain point as the root-cause for everything under a particular node

	//2. getter for a collection of collectors to run next. Each collector represented by name + params bundle

	//3. getter for a collection of diagnosers to run next. Each diagnoser represented by name + params bundle

}
