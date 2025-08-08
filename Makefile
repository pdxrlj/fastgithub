VERSION = 0.0.1

release:
	git tag -d v$(VERSION)
	git push origin --delete v$(VERSION)
	git tag v$(VERSION)
	git push origin v$(VERSION)

# 用户输入指定的版本号，然后执行 release 命令
release-version:
	@read -p "请输入版本号: " VERSION; \
	make release VERSION=$$VERSION