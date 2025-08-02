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

          <!-- 自定义下拉组件 -->
          <div class="custom-select-wrapper">
            <div
              class="custom-select"
              :class="{
                'is-open': isDropdownOpen,
                'is-disabled': isSettingInterval,
              }"
              @click="toggleDropdown"
            >
              <div class="select-display">
                <span class="select-text">{{ getSelectedIntervalText() }}</span>
                <span
                  class="select-arrow"
                  :class="{ 'is-rotated': isDropdownOpen }"
                  >▼</span
                >
              </div>

              <!-- 下拉选项 -->
              <Transition name="dropdown">
                <div
                  v-if="isDropdownOpen"
                  class="dropdown-menu"
                  :class="{ 'dropdown-top': dropdownPosition === 'top' }"
                >
                  <div
                    v-for="option in updateOptions"
                    :key="option.value"
                    class="dropdown-item"
                    :class="{
                      'is-selected': selectedInterval === option.value,
                    }"
                    @click.stop="selectOption(option.value)"
                  >
                    <span class="option-icon">{{ option.icon }}</span>
                    <span class="option-text">{{ option.label }}</span>
                    <span
                      v-if="selectedInterval === option.value"
                      class="option-check"
                      >✓</span
                    >
                  </div>
                </div>
              </Transition>
            </div>
          </div>

          <!-- 更新状态消息 -->
          <Transition name="message">
            <div
              v-if="updateMessage"
              class="update-message"
              :class="updateMessageType"
            >
              {{ updateMessage }}
            </div>
          </Transition>
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
            v-if="activeTab === 'local'"
            class="toolbar-btn"
            @click="openHostsFolder"
            title="打开 Hosts 文件所在文件夹"
          >
            <span class="btn-icon">📁</span>
            打开文件夹
          </button>

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
            <span class="btn-icon">💾</span>
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
                <span v-if="hostsFilePath" class="file-path">
                  📄 {{ hostsFilePath }}
                </span>
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
                        <span class="no-data-icon">📋</span>
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
                      <span class="status-badge active">✅ 活跃</span>
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
                        <span class="no-data-icon">🌍</span>
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
                      <span class="status-badge remote">🌐 远程</span>
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
import { onMounted, ref, computed, onUnmounted } from "vue";
import {
  ReadHosts,
  ReadHostsFromNetwork,
  SetPollInterval,
  WriteHosts,
  GetHostsPath,
  OpenHostsFolder,
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

interface UpdateOption {
  value: number;
  label: string;
  icon: string;
}

const tabs: TabItem[] = [
  { key: "local", label: "本地 Hosts", icon: "📄" },
  { key: "remote", label: "远程 Hosts", icon: "🌍" },
];

// 更新选项配置
const updateOptions: UpdateOption[] = [
  { value: 0, label: "关闭", icon: "⏹️" },
  { value: 3600, label: "每小时", icon: "⏰" },
  { value: 10800, label: "每3小时", icon: "🕐" },
  { value: 21600, label: "每6小时", icon: "🕕" },
  { value: 43200, label: "每12小时", icon: "🕛" },
  { value: 86400, label: "每天", icon: "📅" },
];

const activeTab = ref<"local" | "remote">("local");
const localHosts = ref("");
const remoteHosts = ref("");
const selectedInterval = ref(0);
const isLoading = ref(false);
const hostsFilePath = ref("");
const updateMessage = ref("");
const updateMessageType = ref("success");
const isSettingInterval = ref(false);
const isDropdownOpen = ref(false);
const dropdownPosition = ref<"top" | "bottom">("bottom");
let setPollIntervalTimeout: ReturnType<typeof setTimeout> | null = null;

const parseHosts = (hosts: string): HostEntry[] => {
  if (!hosts) return [];

  const lines = hosts.split("\n");
  const filteredLines = lines
    .map((line) => line.trim())
    .filter((line) => line && !line.startsWith("#"));

  const entries = filteredLines
    .map((line) => {
      const parts = line.split(/\s+/);
      if (parts.length >= 2) {
        return { ip: parts[0], domain: parts.slice(1).join(" ") };
      }
      return null;
    })
    .filter((entry): entry is HostEntry => entry !== null);

  return entries;
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

const getSelectedIntervalText = () => {
  const option = updateOptions.find(
    (opt) => opt.value === selectedInterval.value
  );
  return option ? option.label : "关闭";
};

const refreshCurrentTab = () => {
  if (activeTab.value === "local") {
    loadLocalHosts();
  } else {
    loadRemoteHosts();
  }
};

const openHostsFolder = () => {
  OpenHostsFolder();
};

const loadHostsPath = async () => {
  try {
    hostsFilePath.value = await GetHostsPath();
  } catch (e) {
    console.error("获取hosts文件路径失败:", e);
  }
};

// 下拉组件相关方法
const toggleDropdown = () => {
  if (!isSettingInterval.value) {
    if (!isDropdownOpen.value) {
      // 打开下拉时计算位置
      dropdownPosition.value = getDropdownPosition();
    }
    isDropdownOpen.value = !isDropdownOpen.value;
  }
};

const selectOption = (value: number) => {
  if (selectedInterval.value !== value) {
    selectedInterval.value = value;
    debouncedSetPollInterval();
  }
  isDropdownOpen.value = false;
};

// 点击外部关闭下拉
const handleClickOutside = (event: Event) => {
  const target = event.target as Element;
  if (!target.closest(".custom-select-wrapper")) {
    isDropdownOpen.value = false;
  }
};

// 计算下拉菜单位置
const getDropdownPosition = () => {
  const wrapper = document.querySelector(
    ".custom-select-wrapper"
  ) as HTMLElement;
  if (!wrapper) return "bottom";

  const rect = wrapper.getBoundingClientRect();
  const windowHeight = window.innerHeight;
  const dropdownHeight = Math.min(200, updateOptions.length * 40); // 动态计算高度
  const spaceBelow = windowHeight - rect.bottom;
  const spaceAbove = rect.top;

  // 如果下方空间不足，且上方空间足够，则向上展开
  if (spaceBelow < dropdownHeight && spaceAbove > dropdownHeight) {
    return "top";
  }

  // 如果上下空间都不足，选择空间较大的方向
  if (spaceBelow < dropdownHeight && spaceAbove < dropdownHeight) {
    return spaceAbove > spaceBelow ? "top" : "bottom";
  }

  return "bottom";
};

async function loadLocalHosts() {
  isLoading.value = true;
  try {
    localHosts.value = await ReadHosts();
  } catch (e) {
    console.error("加载本地 Hosts 失败:", e);
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

// 防抖函数，避免频繁调用
function debouncedSetPollInterval() {
  // 清除之前的定时器
  if (setPollIntervalTimeout) {
    clearTimeout(setPollIntervalTimeout);
  }

  // 立即显示正在设置的状态
  isSettingInterval.value = true;
  updateMessage.value = "⏳ 正在设置...";
  updateMessageType.value = "info";

  // 500ms后执行实际的设置操作
  setPollIntervalTimeout = setTimeout(() => {
    setPollInterval();
  }, 500);
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
    updateMessageType.value = "success";

    // 3秒后清除消息
    setTimeout(() => {
      updateMessage.value = "";
    }, 3000);
  } catch (e) {
    // 显示错误消息
    updateMessage.value = `❌ 设置失败：${e}`;
    updateMessageType.value = "error";

    // 5秒后清除消息
    setTimeout(() => {
      updateMessage.value = "";
    }, 5000);
  } finally {
    isSettingInterval.value = false;
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
  loadHostsPath();

  // 添加全局点击事件监听
  document.addEventListener("click", handleClickOutside);
});

onUnmounted(() => {
  // 清理事件监听
  document.removeEventListener("click", handleClickOutside);

  // 清理定时器
  if (setPollIntervalTimeout) {
    clearTimeout(setPollIntervalTimeout);
  }
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

/* 自定义下拉组件样式 */
.custom-select-wrapper {
  position: relative;
  width: 100%;
  z-index: 1000;
}

.custom-select {
  width: 100%;
  background: #34495e;
  border: 1px solid #4a5f7a;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.custom-select:hover:not(.is-disabled) {
  border-color: #5a6f8a;
  background: #3a4f5e;
}

.custom-select.is-open {
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.custom-select.is-disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.select-display {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  color: #ecf0f1;
  font-size: 13px;
}

.select-text {
  font-weight: 500;
}

.select-arrow {
  font-size: 10px;
  transition: transform 0.2s ease;
  color: #bdc3c7;
}

.select-arrow.is-rotated {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #34495e;
  border: 1px solid #4a5f7a;
  border-top: none;
  border-radius: 0 0 6px 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 1000;
  max-height: 200px;
  overflow-y: auto;
  /* 防止超出窗口边界 */
  max-width: 100%;
  box-sizing: border-box;
}

.dropdown-menu.dropdown-top {
  top: auto;
  bottom: 100%;
  border-top: 1px solid #4a5f7a;
  border-bottom: none;
  border-radius: 6px 6px 0 0;
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.3);
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
}

.dropdown-item:hover {
  background: #3a4f5e;
}

.dropdown-item.is-selected {
  background: #3498db;
  color: #ffffff;
}

.option-icon {
  font-size: 14px;
  width: 16px;
  text-align: center;
}

.option-text {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
}

.option-check {
  font-size: 12px;
  color: #27ae60;
  font-weight: bold;
}

.dropdown-item.is-selected .option-check {
  color: #ffffff;
}

/* 动画效果 */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.message-enter-active,
.message-leave-active {
  transition: all 0.3s ease;
}

.message-enter-from,
.message-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}

.update-message {
  margin-top: 8px;
  padding: 6px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  text-align: center;
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

.update-message.info {
  background: rgba(52, 152, 219, 0.2);
  color: #3498db;
  border: 1px solid rgba(52, 152, 219, 0.3);
}

/* 自定义滚动条样式 */
.dropdown-menu::-webkit-scrollbar {
  width: 6px;
}

.dropdown-menu::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.dropdown-menu::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
  transition: background 0.2s ease;
}

.dropdown-menu::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}

/* 表格滚动条美化 */
.table-wrapper::-webkit-scrollbar {
  width: 8px;
}

.table-wrapper::-webkit-scrollbar-track {
  background: #f1f3f4;
  border-radius: 4px;
}

.table-wrapper::-webkit-scrollbar-thumb {
  background: #c1c7cd;
  border-radius: 4px;
  transition: background 0.2s ease;
}

.table-wrapper::-webkit-scrollbar-thumb:hover {
  background: #a8b0b8;
}

/* 主内容区样式优化 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100vh;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 28px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.08);
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
  background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 20px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(5px);
}

.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #27ae60;
  box-shadow: 0 0 0 2px rgba(39, 174, 96, 0.2);
}

.status-indicator.loading .status-dot {
  background: #f39c12;
  animation: pulse 1.5s infinite;
  box-shadow: 0 0 0 2px rgba(243, 156, 18, 0.2);
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.7;
    transform: scale(1.1);
  }
}

.status-text {
  font-size: 13px;
  color: #6c757d;
  font-weight: 500;
}

.toolbar-right {
  display: flex;
  gap: 12px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 18px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  color: #495057;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(5px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.toolbar-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 1);
  border-color: rgba(52, 152, 219, 0.3);
  box-shadow: 0 4px 16px rgba(52, 152, 219, 0.2);
  transform: translateY(-1px);
}

.toolbar-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.toolbar-btn.primary {
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  border-color: #3498db;
  color: #ffffff;
  box-shadow: 0 4px 16px rgba(52, 152, 219, 0.3);
}

.toolbar-btn.primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #2980b9 0%, #1f5f8b 100%);
  border-color: #2980b9;
  box-shadow: 0 6px 20px rgba(52, 152, 219, 0.4);
  transform: translateY(-2px);
}

.btn-icon {
  font-size: 16px;
  display: inline-block;
  transition: transform 0.2s ease;
}

.toolbar-btn:hover .btn-icon {
  transform: scale(1.1);
}

.toolbar-btn.primary:hover .btn-icon {
  transform: scale(1.1) rotate(5deg);
}

/* 内容区域样式优化 */
.content-area {
  flex: 1;
  padding: 28px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.data-container {
  flex: 1;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.hosts-table-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.table-header {
  padding: 20px 24px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 16px 16px 0 0;
}

.table-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.entry-count {
  font-size: 14px;
  color: #6c757d;
  font-weight: 600;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 20px;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.file-path {
  font-size: 12px;
  color: #495057;
  font-weight: 500;
  padding: 4px 8px;
  background: rgba(52, 152, 219, 0.1);
  border-radius: 6px;
  border: 1px solid rgba(52, 152, 219, 0.2);
  font-family: "Consolas", "Monaco", "Courier New", monospace;
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
  border-collapse: separate;
  border-spacing: 0;
  font-size: 14px;
}

.hosts-table th {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  color: #495057;
  font-weight: 700;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 16px 20px;
  text-align: left;
  border-bottom: 2px solid rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1;
  backdrop-filter: blur(5px);
}

.hosts-table td {
  padding: 16px 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  vertical-align: middle;
  text-align: left;
  transition: all 0.2s ease;
}

.table-row:hover {
  background: rgba(52, 152, 219, 0.05);
  transform: translateX(4px);
}

.col-index {
  width: 60px;
  color: #6c757d;
  font-weight: 600;
  text-align: center;
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
  background: linear-gradient(135deg, #f1f3f4 0%, #e8eaed 100%);
  padding: 6px 12px;
  border-radius: 8px;
  font-family: "Consolas", "Monaco", "Courier New", monospace;
  font-size: 13px;
  color: #495057;
  font-weight: 600;
  border: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.domain-text {
  color: #2c3e50;
  font-weight: 600;
  font-size: 14px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
}

.status-badge.active {
  background: linear-gradient(135deg, #d4edda 0%, #c3e6cb 100%);
  color: #155724;
  border: 1px solid rgba(21, 87, 36, 0.2);
}

.status-badge.remote {
  background: linear-gradient(135deg, #cce5ff 0%, #b3d9ff 100%);
  color: #004085;
  border: 1px solid rgba(0, 64, 133, 0.2);
}

.status-badge:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.no-data {
  text-align: center;
  padding: 60px 20px;
}

.no-data-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.no-data-icon {
  font-size: 48px;
  opacity: 0.6;
  filter: grayscale(0.3);
}

.no-data-text {
  color: #6c757d;
  font-size: 16px;
  font-weight: 500;
}

/* 响应式设计优化 */
@media (max-width: 1024px) {
  .sidebar {
    width: 220px;
  }

  .content-area {
    padding: 20px;
  }

  .toolbar {
    padding: 16px 20px;
  }

  .page-title {
    font-size: 20px;
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

  .toolbar {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .toolbar-left {
    justify-content: center;
  }

  .toolbar-right {
    justify-content: center;
  }
}
</style>
