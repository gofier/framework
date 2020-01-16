package policy

type RoutePolicier interface {
	Can(policy IPolicy, action Action)
}
