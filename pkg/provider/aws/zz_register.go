package aws

func (p *Provider) registerGeneratedTypes(mapping map[string]mapper) {
	p.register_ec2(mapping)
	p.register_elb(mapping)
	p.register_lambda(mapping)
	p.register_rds(mapping)
	p.register_s3(mapping)
}
