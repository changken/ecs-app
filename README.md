# ecs-app

ECS Fargate 示範應用程式，搭配 CodeDeploy **Blue/Green** 自動部署。

基礎設施定義（Terraform）位於 [changken/infra-lab](https://github.com/changken/infra-lab)
的 `terraform/envs/aws-ecs-fargate/`。

## 架構

```
push to main
    │
    ▼
GitHub Actions
    ├── docker build → ECR push（:$SHA + :latest）
    ├── 更新 ECS Task Definition（換 image）
    ├── CodeDeploy create-deployment
    │       ├── Green TG 啟動新 tasks
    │       ├── :8080 Test Listener 可預覽新版本
    │       └── AllAtOnce → :80 Production Listener 切換
    └── smoke test curl /health
```

## Endpoints

| Path | 說明 |
|------|------|
| `GET /` | JSON：status、version、hostname、ECS task metadata |
| `GET /health` | plaintext `ok`（ALB health check 用）|
| `GET /version` | JSON：`version`、`git_commit`（部署驗證用）|

## 環境變數

| 變數 | 來源 | 說明 |
|------|------|------|
| `APP_VERSION` | ECS Task Definition（Terraform）| 應用版本號 |
| `GIT_COMMIT` | Docker build-arg（GitHub Actions）| 7 碼 commit SHA |
| `AWS_REGION` | ECS 自動注入 | 所在 region |
| `PORT` | 選填，預設 `8080` | 監聽 port |

## 本地開發

```bash
# 需要 Go 1.23+
go mod init app
go run main.go
# → http://localhost:8080
```

```bash
# Docker
docker build -t ecs-app .
docker run -p 8080:8080 -e APP_VERSION=local ecs-app
curl localhost:8080/version
```

## CI/CD 流程

push 到 `main` 自動觸發 `.github/workflows/deploy.yml`：

```
1. Configure AWS credentials（OIDC，不需長效 Access Key）
2. ECR login
3. docker build --build-arg GIT_COMMIT=<sha7> → push :$SHA + :latest
4. aws ecs describe-task-definition → jq 換 image → register 新版 task def
5. aws deploy create-deployment（CodeDeploy Blue/Green）
6. aws deploy wait deployment-successful
7. curl /health smoke test
```

### 所需 GitHub Secret / Variable

無需設定任何 Secret。OIDC role ARN 直接寫在 workflow，
由 `infra-lab` Terraform 建立的 `infra-lab-dev-github-actions-role` 授權。

## 部署驗證

```bash
ALB="http://infra-lab-dev-alb-17053701.us-east-1.elb.amazonaws.com"

# 部署期間：透過 Test Listener（:8080）看 Green TG 新版本
curl -s $ALB:8080/version

# 部署完成後：Production Listener（:80）已切換
curl -s $ALB/version
# {"version":"1.0.0","git_commit":"d3e8f97"}
```

## 手動觸發

```bash
# GitHub CLI
gh workflow run deploy.yml --repo changken/ecs-app
```

或 GitHub UI → Actions → Deploy to ECS Fargate (Blue/Green) → Run workflow
