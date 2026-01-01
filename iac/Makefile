.PHONY: init-cf update-cf-defaults help

# CloudFormationテンプレートのパス
CF_TEMPLATE := cloudformations/github-actions-role.yaml

# Gitリモートから組織名とリポジトリ名を抽出
REMOTE_URL := $(shell git remote get-url origin 2>/dev/null)
# HTTPS形式: https://github.com/org/repo.git
# SSH形式: git@github.com:org/repo.git
GH_ORG := $(shell echo "$(REMOTE_URL)" | sed -E 's|.*github\.com[:/]([^/]+)/.*|\1|')
GH_REPO := $(shell echo "$(REMOTE_URL)" | sed -E 's|.*github\.com[:/][^/]+/([^.]+)(\.git)?|\1|')

# CloudFormationテンプレートのデフォルト値を更新
update-cf-defaults:
	@echo "=== Updating CloudFormation Template Defaults ==="
	@echo "Git Remote URL: $(REMOTE_URL)"
	@echo "GitHub Organization: $(GH_ORG)"
	@echo "GitHub Repository: $(GH_REPO)"
	@echo ""
	@if [ -z "$(GH_ORG)" ] || [ -z "$(GH_REPO)" ]; then \
		echo "❌ Error: Could not extract GitHub org/repo from git remote"; \
		exit 1; \
	fi
	@sed -i '' "s|Default: 'YOUR_GITHUB_USERNAME'|Default: '$(GH_ORG)'|g" $(CF_TEMPLATE)
	@sed -i '' "s|Default: 'automation-deploy-template'|Default: '$(GH_REPO)'|g" $(CF_TEMPLATE)
	@echo "✅ Updated $(CF_TEMPLATE)"
	@echo ""
	@echo "Updated values:"
	@grep -A1 "GitHubOrganization:" $(CF_TEMPLATE) | grep "Default:" || true
	@grep -A1 "GitHubRepository:" $(CF_TEMPLATE) | grep "Default:" || true

# CloudFormationスタックの初期化（デフォルト値更新 + デプロイ案内）
init-cf: update-cf-defaults
	@echo ""
	@echo "=== Next Steps ==="
	@echo ""
	@echo "1. Deploy the CloudFormation stack:"
	@echo "   aws cloudformation create-stack \\"
	@echo "     --stack-name $(GH_REPO)-github-actions \\"
	@echo "     --template-body file://$(CF_TEMPLATE) \\"
	@echo "     --capabilities CAPABILITY_NAMED_IAM \\"
	@echo "     --parameters \\"
	@echo "       ParameterKey=RoleName,ParameterValue=$(GH_REPO)-GitHubActionsRole"
	@echo ""
	@echo "2. After deployment, add these secrets to GitHub repository:"
	@echo "   - AWS_ROLE_TO_ASSUME: (from CloudFormation output)"
	@echo "   - AWS_REGION: ap-northeast-1"
	@echo "   - PROJECT_NAME: $(GH_REPO)"
	@echo "   - ROOT_DOMAIN: your-domain.com"

# ヘルプ
help:
	@echo "使用可能なコマンド:"
	@echo "  make update-cf-defaults - GitリモートからCloudFormationテンプレートのデフォルト値を更新"
	@echo "  make init-cf            - デフォルト値を更新してデプロイ手順を表示"
	@echo "  make help               - このヘルプを表示"
