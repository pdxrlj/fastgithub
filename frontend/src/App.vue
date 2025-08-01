<script lang="ts" setup>
import hosts from "components/hosts.vue";
import {
  WindowMinimise,
  WindowToggleMaximise,
  Quit,
} from "../wailsjs/runtime/runtime";

const close = () => Quit();
const mini = () => WindowMinimise(); // 最小化
const maxi = () => WindowToggleMaximise(); // 最大化/还原
</script>

<template>
  <div class="app" style="--wails-draggable: drag">
    <div class="header">
      <div class="title">GitHub Fast</div>
      <div class="controls">
        <button @click="mini" class="control-btn minimize" title="最小化">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
            <rect x="4" y="11" width="16" height="2" rx="1" />
          </svg>
        </button>
        <button @click="maxi" class="control-btn maximize" title="最大化">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
            <path d="M4 4h16v16H4V4zm2 2v12h12V6H6z" />
          </svg>
        </button>
        <button @click="close" class="control-btn close" title="关闭">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
            <path
              d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
            />
          </svg>
        </button>
      </div>
    </div>
    <div class="content">
      <hosts></hosts>
    </div>
  </div>
</template>

<style>
.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #1e1e1e;
  color: #ffffff;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC",
    "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 32px;
  background: linear-gradient(135deg, #2d2d2d 0%, #1a1a1a 100%);
  border-bottom: 1px solid #333;
  -webkit-app-region: drag;
  padding: 0 12px;
}

.title {
  font-size: 13px;
  font-weight: 500;
  color: #e0e0e0;
  letter-spacing: 0.5px;
}

.controls {
  display: flex;
  align-items: center;
  gap: 2px;
  -webkit-app-region: no-drag;
}

.control-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 46px;
  height: 32px;
  border: none;
  background: transparent;
  color: #b0b0b0;
  cursor: pointer;
  transition: all 0.2s ease;
  border-radius: 4px;
  position: relative;
  overflow: hidden;
}

.control-btn:hover {
  color: #ffffff;
  background: linear-gradient(135deg, #3a3a3a 0%, #2a2a2a 100%);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.control-btn:active {
  transform: scale(0.95);
  background: linear-gradient(135deg, #404040 0%, #303030 100%);
}

.control-btn.close:hover {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
  color: #ffffff;
}

.control-btn svg {
  transition: transform 0.2s ease;
}

.control-btn:hover svg {
  transform: scale(1.1);
}

.control-btn:active svg {
  transform: scale(0.9);
}

.content {
  flex: 1;
  overflow: hidden;
  background: #1e1e1e;
}
</style>
