<template>
  <div class="hosts-app">
    <!-- 侧边栏 -->
    <div class="sidebar">
      <div class="sidebar-header">
        <div class="app-logo">
          <img src="../assets/logo.svg" alt="GitHub Fast" class="logo-icon" />
          <span class="logo-text">GitHub Fast</span>
        </div>
      </div>

      <nav class="sidebar-nav">
        <div
          v-for="tab in tabs"
          :key="tab.key"
          :class="['nav-item', { active: activeTab === tab.key }]"
          @click="activeTab = tab.key"
        >
          <span class="nav-icon">{{ tab.icon }}</span>
          <span class="nav-label">{{ tab.label }}</span>
        </div>
      </nav>

      <div class="sidebar-footer">
        <div class="auto-update-section">
          <div class="section-title">自动更新</div>
          <select
            v-model="selectedInterval"
            @change="setPollInterval"
            class="update-select"
          >
            <option :value="0">关闭</option>
            <option :value="3600">每小时</option>
            <option :value="10800">每3小时</option>
            <option :value="21600">每6小时</option>
            <option :value="43200">每12小时</option>
            <option :value="86400">每天</option>
          </select>
          <div v-if="updateMessage" class="update-message" :class="updateMessageType">
            {{ updateMessage }}
          </div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 工具栏 -->
      <div class="toolbar">
        <div class="toolbar-left">
          <h1 class="page-title">{{ getCurrentTabLabel() }}</h1>
          <div class="status-indicator" :class="{ loading: isLoading }">
            <span class="status-dot"></span>
            <span class="status-text">{{ getStatusText() }}</span>
          </div>
        </div>

        <div class="toolbar-right">
          <button
            class="toolbar-btn"
            @click="refreshCurrentTab"
            :disabled="isLoading"
          >
            <span class="btn-icon">🔄</span>
            刷新
          </button>

          <button
            v-if="activeTab === 'remote'"
            class="toolbar-btn primary"
            @click="applyRemoteHosts"
            :disabled="!canApplyRemote"
          >
            <span class="btn-icon">✅</span>
            应用到本地
          </button>
        </div>
      </div>

      <!-- 数据展示区 -->
      <div class="content-area">
        <div class="data-container">
          <!-- 本地 Hosts -->
          <div v-if="activeTab === 'local'" class="hosts-table-container">
            <div class="table-header">
              <div class="table-info">
                <span class="entry-count"
                  >共 {{ parsedLocalHosts.length }} 条记录</span
                >
              </div>
            </div>

            <div class="table-wrapper">
              <table class="hosts-table">
                <thead>
                  <tr>
                    <th class="col-index">#</th>
                    <th class="col-ip">IP 地址</th>
                    <th class="col-domain">域名</th>
                    <th class="col-status">状态</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-if="parsedLocalHosts.length === 0">
                    <td colspan="4" class="no-data">
                      <div class="no-data-content">
                        <span class="no-data-icon">📝</span>
                        <span class="no-data-text">{{
                          localHosts ? "暂无有效的 Hosts 记录" : "正在加载..."
                        }}</span>
                      </div>
                    </td>
                  </tr>
                  <tr
                    v-for="(host, index) in parsedLocalHosts"
                    :key="index"
                    class="table-row"
                  >
                    <td class="col-index">{{ index + 1 }}</td>
                    <td class="col-ip">
                      <code class="ip-code">{{ host.ip }}</code>
                    </td>
                    <td class="col-domain">
                      <span class="domain-text">{{ host.domain }}</span>
                    </td>
                    <td class="col-status">
                      <span class="status-badge active">活跃</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- 远程 Hosts -->
          <div v-if="activeTab === 'remote'" class="hosts-table-container">
            <div class="table-header">
              <div class="table-info">
                <span class="entry-count"
                  >共 {{ parsedRemoteHosts.length }} 条记录</span
                >
              </div>
            </div>

            <div class="table-wrapper">
              <table class="hosts-table">
                <thead>
                  <tr>
                    <th class="col-index">#</th>
                    <th class="col-ip">IP 地址</th>
                    <th class="col-domain">域名</th>
                    <th class="col-status">状态</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-if="parsedRemoteHosts.length === 0">
                    <td colspan="4" class="no-data">
                      <div class="no-data-content">
                        <span class="no-data-icon">🌐</span>
                        <span class="no-data-text">{{
                          remoteHosts ? "暂无有效的 Hosts 记录" : "正在加载..."
                        }}</span>
                      </div>
                    </td>
                  </tr>
                  <tr
                    v-for="(host, index) in parsedRemoteHosts"
                    :key="index"
                    class="table-row"
                  >
                    <td class="col-index">{{ index + 1 }}</td>
                    <td class="col-ip">
                      <code class="ip-code">{{ host.ip }}</code>
                    </td>
                    <td class="col-domain">
                      <span class="domain-text">{{ host.domain }}</span>
                    </td>
                    <td class="col-status">
                      <span class="status-badge remote">远程</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, computed } from "vue";
import {
  ReadHosts,
  ReadHostsFromNetwork,
  SetPollInterval,
  WriteHosts,
} from "../../wailsjs/go/main/Host";

interface HostEntry {
  ip: string;
  domain: string;
}

interface TabItem {
  key: "local" | "remote";
  label: string;
  icon: string;
}

const tabs: TabItem[] = [
  { key: "local", label: "本地 Hosts", icon: "💻" },
  { key: "remote", label: "远程 Hosts", icon: "🌐" },
];

const activeTab = ref<"local" | "remote">("local");
const localHosts = ref("");
const remoteHosts = ref("");
const selectedInterval = ref(0)
const isLoading = ref(false)
const updateMessage = ref('')
const updateMessageType = ref('success');

const parseHosts = (hosts: string): HostEntry[] => {
  if (!hosts) return [];
  return hosts
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line && !line.startsWith("#"))
    .map((line) => {
      const parts = line.split(/\s+/);
      if (parts.length >= 2) {
        return { ip: parts[0], domain: parts.slice(1).join(" ") };
      }
      return null;
    })
    .filter((entry): entry is HostEntry => entry !== null);
};

const parsedLocalHosts = computed(() => parseHosts(localHosts.value));
const parsedRemoteHosts = computed(() => parseHosts(remoteHosts.value));

const canApplyRemote = computed(() => {
  return (
    remoteHosts.value &&
    !remoteHosts.value.startsWith("❌") &&
    parsedRemoteHosts.value.length > 0
  );
});

const getCurrentTabLabel = () => {
  const tab = tabs.find((t) => t.key === activeTab.value);
  return tab ? tab.label : "";
};

const getStatusText = () => {
  if (isLoading.value) return "加载中...";
  if (activeTab.value === "local") {
    return parsedLocalHosts.value.length > 0 ? "已加载" : "无数据";
  } else {
    return parsedRemoteHosts.value.length > 0 ? "已加载" : "无数据";
  }
};

const refreshCurrentTab = () => {
  if (activeTab.value === "local") {
    loadLocalHosts();
  } else {
    loadRemoteHosts();
  }
};

async function loadLocalHosts() {
  isLoading.value = true;
  try {
    localHosts.value = await ReadHosts();
  } catch (e) {
    localHosts.value = `❌ 读取失败：${e}`;
  } finally {
    isLoading.value = false;
  }
}

async function loadRemoteHosts() {
  isLoading.value = true;
  try {
    remoteHosts.value = await ReadHostsFromNetwork();
  } catch (e) {
    remoteHosts.value = `❌ 获取失败：${e}`;
  } finally {
    isLoading.value = false;
  }
}

async function setPollInterval() {
  try {
    await SetPollInterval(selectedInterval.value);
    const msg =
      selectedInterval.value === 0
        ? "已关闭自动更新"
        : `已设置${selectedInterval.value / 3600}小时自动更新`;
    
    // 显示成功消息
    updateMessage.value = `✅ ${msg}`;
    updateMessageType.value = 'success';
    
    // 3秒后清除消息
    setTimeout(() => {
      updateMessage.value = '';
    }, 3000);
  } catch (e) {
    // 显示错误消息
    updateMessage.value = `❌ 设置失败：${e}`;
    updateMessageType.value = 'error';
    
    // 5秒后清除消息
    setTimeout(() => {
      updateMessage.value = '';
    }, 5000);
  }
}

async function applyRemoteHosts() {
  if (!canApplyRemote.value) {
    console.warn("❌ 远程 Hosts 无效，无法应用");
    return;
  }

  isLoading.value = true;
  try {
    await WriteHosts(remoteHosts.value);
    await loadLocalHosts();
    console.log("✅ 已成功应用远程 Hosts！");
  } catch (e) {
    console.error(`❌ 应用失败：${e}`);
  } finally {
    isLoading.value = false;
  }
}

onMounted(() => {
  loadLocalHosts();
  loadRemoteHosts();
});
</script>

<style scoped>
.hosts-app {
  display: flex;
  height: 100vh;
  background: #f8f9fa;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC",
    "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
  overflow: hidden;
}

/* 侧边栏样式 */
.sidebar {
  width: 260px;
  background: #2c3e50;
  color: #ecf0f1;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #34495e;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #34495e;
}

.app-logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #ecf0f1;
}

.sidebar-nav {
  flex: 1;
  padding: 20px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  border-left: 3px solid transparent;
}

.nav-item:hover {
  background: #34495e;
}

.nav-item.active {
  background: #3498db;
  border-left-color: #2980b9;
}

.nav-icon {
  font-size: 16px;
  width: 20px;
  text-align: center;
}

.nav-label {
  font-size: 14px;
  font-weight: 500;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid #34495e;
  min-height: 120px;
  display: flex;
  flex-direction: column;
}

.auto-update-section {
  margin-bottom: 0;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  color: #bdc3c7;
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.update-select {
  width: 100%;
  padding: 8px 12px;
  background: #34495e;
  border: 1px solid #4a5f7a;
  border-radius: 6px;
  color: #ecf0f1;
  font-size: 13px;
  cursor: pointer;
}

.update-select:focus {
  outline: none;
  border-color: #3498db;
}

.update-message {
  margin-top: 8px;
  padding: 6px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  text-align: center;
  animation: fadeIn 0.3s ease;
}

.update-message.success {
  background: rgba(39, 174, 96, 0.2);
  color: #27ae60;
  border: 1px solid rgba(39, 174, 96, 0.3);
}

.update-message.error {
  background: rgba(231, 76, 60, 0.2);
  color: #e74c3c;
  border: 1px solid rgba(231, 76, 60, 0.3);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 主内容区样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100vh;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: #ffffff;
  border-bottom: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #27ae60;
}

.status-indicator.loading .status-dot {
  background: #f39c12;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.status-text {
  font-size: 13px;
  color: #7f8c8d;
}

.toolbar-right {
  display: flex;
  gap: 12px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  background: #ffffff;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  color: #495057;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.toolbar-btn:hover:not(:disabled) {
  background: #f8f9fa;
  border-color: #adb5bd;
}

.toolbar-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.toolbar-btn.primary {
  background: #3498db;
  border-color: #3498db;
  color: #ffffff;
}

.toolbar-btn.primary:hover:not(:disabled) {
  background: #2980b9;
  border-color: #2980b9;
}

.btn-icon {
  font-size: 14px;
}

/* 内容区域样式 */
.content-area {
  flex: 1;
  padding: 24px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.data-container {
  flex: 1;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.hosts-table-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.table-header {
  padding: 16px 20px;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.table-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.entry-count {
  font-size: 13px;
  color: #6c757d;
  font-weight: 500;
}

.table-wrapper {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 0;
  padding-bottom: 20px;
}

.hosts-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.hosts-table th {
  background: #f8f9fa;
  color: #495057;
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 12px 16px;
  text-align: left;
  border-bottom: 2px solid #e9ecef;
  position: sticky;
  top: 0;
  z-index: 1;
}

.hosts-table td {
  padding: 12px 16px;
  border-bottom: 1px solid #f1f3f4;
  vertical-align: middle;
}

.table-row:hover {
  background: #f8f9fa;
}

.col-index {
  width: 60px;
  text-align: center;
  color: #6c757d;
  font-weight: 500;
}

.col-ip {
  width: 140px;
}

.col-domain {
  flex: 1;
}

.col-status {
  width: 80px;
}

.ip-code {
  background: #f1f3f4;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: "Consolas", "Monaco", monospace;
  font-size: 13px;
  color: #495057;
}

.domain-text {
  color: #2c3e50;
  font-weight: 500;
}

.status-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-badge.active {
  background: #d4edda;
  color: #155724;
}

.status-badge.remote {
  background: #cce5ff;
  color: #004085;
}

.no-data {
  text-align: center;
  padding: 40px 20px;
}

.no-data-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.no-data-icon {
  font-size: 32px;
  opacity: 0.5;
}

.no-data-text {
  color: #6c757d;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .sidebar {
    width: 220px;
  }

  .content-area {
    padding: 16px;
  }
}

@media (max-width: 768px) {
  .hosts-app {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    height: auto;
    flex-direction: row;
  }

  .sidebar-header {
    flex: 1;
  }

  .sidebar-nav {
    flex: 2;
    display: flex;
    padding: 0;
  }

  .sidebar-footer {
    flex: 1;
  }

  .nav-item {
    flex: 1;
    justify-content: center;
    border-left: none;
    border-bottom: 3px solid transparent;
  }

  .nav-item.active {
    border-bottom-color: #2980b9;
  }
}
</style>
