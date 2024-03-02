# High-Performance Load Balancer

## Conceptual Architecture

- Global Load Balancer: Google's Global Load Balancer will be the frontend, distributing traffic across regional load balancers based on factors like user location and backend health.
- Regional Go Load Balancers: In each region, you'll have custom Go-based load balancers performing health checks and routing traffic to Managed Instance Groups.
- Managed Instance Groups (MIGs): These will contain the instances of your application, scaling up or down as needed.