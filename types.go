/**
 * @package go-response (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Feb 18, 2026; 9:36 PM
 */

package response

type mimetypes struct {
	JSON       string
	XML        string
	HTML       string
	CSS        string
	JAVASCRIPT string
	FILE       string
}

type seconds struct {
	SixHours int
	HalfDay  int
	OneDay   int
	OneWeek  int
	OneMonth int
}

type status struct {
	Ok                    int
	NoContent             int
	BadRequest            int
	Unauthorized          int
	Forbidden             int
	NotFound              int
	MethodNotAllowed      int
	Conflict              int
	ServerError           int
	ServerBadGateway      int
	ServiceGatewayTimeout int
}

type redirect struct {
	Permanent int
	Temporary int
}
