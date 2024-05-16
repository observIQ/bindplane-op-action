package model

type Kind string

const (
	KindProfile                    Kind = "Profile"
	KindContext                    Kind = "Context"
	KindConfiguration              Kind = "Configuration"
	KindAgent                      Kind = "Agent"
	KindAgentVersion               Kind = "AgentVersion"
	KindSource                     Kind = "Source"
	KindProcessor                  Kind = "Processor"
	KindDestination                Kind = "Destination"
	KindExtension                  Kind = "Extension"
	KindSourceType                 Kind = "SourceType"
	KindProcessorType              Kind = "ProcessorType"
	KindDestinationType            Kind = "DestinationType"
	KindExtensionType              Kind = "ExtensionType"
	KindRecommendationType         Kind = "RecommendationType"
	KindUnknown                    Kind = "Unknown"
	KindRollout                    Kind = "Rollout"
	KindRolloutType                Kind = "RolloutType"
	KindOrganization               Kind = "Organization"
	KindAccount                    Kind = "Account"
	KindInvitation                 Kind = "Invitation"
	KindLogin                      Kind = "Login"
	KindUser                       Kind = "User"
	KindAccountOrganizationBinding Kind = "AccountOrganizationBinding"
	KindUserOrganizationBinding    Kind = "UserOrganizationBinding"
	KindUserAccountBinding         Kind = "UserAccountBinding"
	KindAuditEvent                 Kind = "AuidtTrail"
)
