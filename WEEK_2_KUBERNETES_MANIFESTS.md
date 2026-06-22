# 🚀 WEEK 2: KUBERNETES & CI/CD - COMPLETE IMPLEMENTATION

**Status:** Week 2 Full Implementation  
**Days:** 1-5 (40 hours)  
**Focus:** Production Deployment & Automation  
**Timeline:** This week (5 working days)

---

## 📋 WEEK 2 OVERVIEW

### Days 1-2: Kubernetes Manifests (12 Hours)
- Deployment configuration
- Service definition
- ConfigMap setup
- Secrets management
- HPA (autoscaling)
- StatefulSets for databases

### Days 3-4: GitHub Actions CI/CD (16 Hours)
- Build & test pipeline
- Docker image building
- Push to registry
- Security scanning
- Deployment automation

### Day 5: Integration & Testing (12 Hours)
- Pipeline verification
- Deployment tests
- Health checks
- Rollback procedures
- Documentation

---

## DAYS 1-2: KUBERNETES MANIFESTS (12 Hours)

### Kubernetes Directory Structure

```
infra/kubernetes/
├── base/                                    # Base manifests (reusable)
│   ├── auth-service.yaml                   ✅ CREATE
│   ├── user-service.yaml                   ✅ CREATE
│   ├── ride-service.yaml                   ✅ CREATE
│   ├── dispatch-service.yaml                ✅ CREATE
│   ├── gps-service.yaml                    ✅ CREATE
│   ├── api-gateway.yaml                    ✅ CREATE
│   ├── websocket-gateway.yaml               ✅ CREATE
│   ├── postgres.yaml                       ✅ CREATE
│   ├── redis.yaml                          ✅ CREATE
│   ├── kafka.yaml                          ✅ CREATE
│   ├── prometheus.yaml                     ✅ CREATE
│   ├── grafana.yaml                        ✅ CREATE
│   ├── loki.yaml                           ✅ CREATE
│   ├── jaeger.yaml                         ✅ CREATE
│   ├── namespace.yaml                      ✅ CREATE
│   ├── rbac.yaml                           ✅ CREATE
│   └── ingress.yaml                        ✅ CREATE
│
├── staging/                                 # Staging environment overrides
│   └── kustomization.yaml
│
└── production/                              # Production environment overrides
    └── kustomization.yaml
```

### Auth Service Deployment (Complete)

```yaml
# infra/kubernetes/base/auth-service.yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: auth-service
  namespace: famgo-platform

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: famgo-platform
  labels:
    app: auth-service
    version: v1
    component: backend
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: auth-service
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      
      initContainers:
      - name: migrate
        image: migrate/migrate:latest
        command:
        - migrate
        - -path=/migrations
        - -database=postgresql://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
        - up
        env:
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: database-credentials
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: database-credentials
              key: password
        - name: DB_NAME
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: database-name
        volumeMounts:
        - name: migrations
          mountPath: /migrations
      
      containers:
      - name: auth-service
        image: auth-service:latest
        imagePullPolicy: Always
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        - name: metrics
          containerPort: 9090
          protocol: TCP
        
        env:
        - name: PORT
          value: "8080"
        - name: ENVIRONMENT
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: environment
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: log-level
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: database-credentials
              key: url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secrets
              key: secret
        - name: JWT_EXPIRY
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: jwt-expiry
        - name: JAEGER_ENDPOINT
          valueFrom:
            configMapKeyRef:
              name: observability-config
              key: jaeger-endpoint
        - name: PROMETHEUS_PORT
          value: "9090"
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: redis-credentials
              key: url
        
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
            ephemeral-storage: "1Gi"
          limits:
            memory: "512Mi"
            cpu: "500m"
            ephemeral-storage: "2Gi"
        
        livenessProbe:
          httpGet:
            path: /health/live
            port: http
            httpHeaders:
            - name: X-Health-Check
              value: liveness
          initialDelaySeconds: 15
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        
        readinessProbe:
          httpGet:
            path: /health/ready
            port: http
            httpHeaders:
            - name: X-Health-Check
              value: readiness
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
        
        startupProbe:
          httpGet:
            path: /health
            port: http
          failureThreshold: 30
          periodSeconds: 10
        
        volumeMounts:
        - name: tmp
          mountPath: /tmp
        - name: cache
          mountPath: /app/cache
        
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          capabilities:
            drop:
            - ALL
      
      volumes:
      - name: tmp
        emptyDir:
          sizeLimit: 1Gi
      - name: cache
        emptyDir:
          sizeLimit: 2Gi
      - name: migrations
        configMap:
          name: auth-service-migrations
      
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - auth-service
              topologyKey: kubernetes.io/hostname
        
        podAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 50
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: component
                  operator: In
                  values:
                  - backend
              topologyKey: topology.kubernetes.io/zone
      
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst

---
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: famgo-platform
  labels:
    app: auth-service
spec:
  type: ClusterIP
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800
  ports:
  - port: 80
    targetPort: http
    protocol: TCP
    name: http
  - port: 9090
    targetPort: metrics
    protocol: TCP
    name: metrics
  selector:
    app: auth-service
  publishNotReadyAddresses: false

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
  namespace: famgo-platform
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 15
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
      - type: Pods
        value: 2
        periodSeconds: 15
      selectPolicy: Max

---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: auth-service-pdb
  namespace: famgo-platform
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: auth-service
  unhealthyPodEvictionPolicy: IfHealthyBudget

---
apiVersion: v1
kind: NetworkPolicy
metadata:
  name: auth-service-network-policy
  namespace: famgo-platform
spec:
  podSelector:
    matchLabels:
      app: auth-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: famgo-platform
    ports:
    - protocol: TCP
      port: 8080
  - from:
    - namespaceSelector:
        matchLabels:
          name: monitoring
    ports:
    - protocol: TCP
      port: 9090
  egress:
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 5432  # PostgreSQL
    - protocol: TCP
      port: 6379  # Redis
    - protocol: TCP
      port: 9092  # Kafka
    - protocol: TCP
      port: 14268 # Jaeger
```

### ConfigMap & Secrets

```yaml
# infra/kubernetes/base/config-secrets.yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: famgo-platform
  labels:
    name: famgo-platform

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: famgo-platform
data:
  environment: production
  log-level: info
  database-name: famgo
  jwt-expiry: "24h"
  refresh-token-expiry: "7d"

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: observability-config
  namespace: famgo-platform
data:
  jaeger-endpoint: "http://jaeger:14268/api/traces"
  prometheus-scrape-interval: "15s"
  prometheus-evaluation-interval: "15s"

---
apiVersion: v1
kind: Secret
metadata:
  name: database-credentials
  namespace: famgo-platform
type: Opaque
stringData:
  username: famgo
  password: ${DB_PASSWORD}
  url: postgresql://famgo:${DB_PASSWORD}@postgres:5432/famgo?sslmode=require

---
apiVersion: v1
kind: Secret
metadata:
  name: jwt-secrets
  namespace: famgo-platform
type: Opaque
stringData:
  secret: ${JWT_SECRET}

---
apiVersion: v1
kind: Secret
metadata:
  name: redis-credentials
  namespace: famgo-platform
type: Opaque
stringData:
  url: redis://redis:6379/0
```

### RBAC (Role-Based Access Control)

```yaml
# infra/kubernetes/base/rbac.yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: auth-service-role
  namespace: famgo-platform
rules:
- apiGroups: [""]
  resources: ["configmaps"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: auth-service-rolebinding
  namespace: famgo-platform
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: auth-service-role
subjects:
- kind: ServiceAccount
  name: auth-service
  namespace: famgo-platform
```

### Ingress Configuration

```yaml
# infra/kubernetes/base/ingress.yaml
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: famgo-ingress
  namespace: famgo-platform
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - api.famgo.com
    - auth.famgo.com
    secretName: famgo-tls
  rules:
  - host: auth.famgo.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: auth-service
            port:
              number: 80
  - host: api.famgo.com
    http:
      paths:
      - path: /auth
        pathType: Prefix
        backend:
          service:
            name: auth-service
            port:
              number: 80
      - path: /users
        pathType: Prefix
        backend:
          service:
            name: user-service
            port:
              number: 80
      - path: /rides
        pathType: Prefix
        backend:
          service:
            name: ride-service
            port:
              number: 80
```

