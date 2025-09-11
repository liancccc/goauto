// GOAUTO 前端交互脚本

// 全局变量
let serverStats = {
    hostname: '加载中...',
    cpuUsage: 0,
    memoryUsage: 0,
    memoryTotal: 0
};

// 页面加载完成后初始化 - 将在页面检测逻辑中处理

// 初始化应用
function initializeApp() {
    loadServerStats();
    loadTaskList();
    setupEventListeners();
    
    // 定期刷新任务列表
    setInterval(loadTaskList, 10000); // 每10秒刷新一次
}

// 加载服务器统计信息
function loadServerStats() {
    // 从API获取服务器信息
    fetchSystemInfo();
    
    // 定期更新服务器状态
    setInterval(fetchSystemInfo, 5000);
}

// 从API获取系统信息
function fetchSystemInfo() {
    fetch('/system/info')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                // 处理系统信息，即使某些字段为空也是正常的
                const systemInfo = data.data || {};
                const stats = {
                    hostname: systemInfo.hostname || 'localhost',
                    cpuUsage: Math.round(systemInfo.cpu_usage || 0),
                    memoryUsage: systemInfo.memory_used || 0,
                    memoryTotal: systemInfo.memory_total || 8
                };
                updateServerStats(stats);
                console.log('系统信息更新成功:', stats);
            } else {
                throw new Error(data.message || '获取系统信息失败');
            }
        })
        .catch(error => {
            console.error('获取系统信息失败:', error);
            showNotification(`获取系统信息失败: ${error.message}`, 'error');
            // 如果API调用失败，使用默认值
            updateServerStats(serverStats);
        });
}

// 更新服务器统计信息
function updateServerStats(stats) {
    const hostnameEl = document.getElementById('hostname');
    const cpuUsageEl = document.getElementById('cpu-usage');
    const memoryUsageEl = document.getElementById('memory-usage');
    
    if (hostnameEl) hostnameEl.textContent = stats.hostname;
    if (cpuUsageEl) {
        cpuUsageEl.textContent = stats.cpuUsage + '%';
        cpuUsageEl.className = 'stat-value text-base-content';
    }
    if (memoryUsageEl) {
        const memoryUsedGB = stats.memoryUsage.toFixed(1);
        const memoryTotalGB = stats.memoryTotal.toFixed(1);
        memoryUsageEl.textContent = memoryUsedGB + 'GB';
        memoryUsageEl.className = 'stat-value text-base-content';
        
        // 更新内存描述信息
        const memoryDescEl = document.getElementById('memory-desc');
        if (memoryDescEl) {
            memoryDescEl.textContent = `已使用 / ${memoryTotalGB}GB`;
        }
    }
}

// 加载任务列表
function loadTaskList() {
    // 从API获取任务列表
    fetch('/task/list')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                // 处理任务列表，即使为空数组也是正常的
                const tasks = data.data ? data.data.map(runtime => ({
                    id: runtime.task,
                    name: runtime.task,
                    mode: runtime.flow, // 修复：使用 flow 而不是 mode
                    status: runtime.status,
                    startTime: formatTime(runtime.start_at),
                    endTime: runtime.end_at ? formatTime(runtime.end_at) : null,
                    description: '安全扫描任务',
                    command: runtime.command,
                    pid: runtime.pid
                })) : [];
                updateTaskList(tasks);
                console.log('任务列表更新成功:', tasks);
            } else {
                throw new Error(data.message || '获取任务列表失败');
            }
        })
        .catch(error => {
            console.error('获取任务列表失败:', error);
            showNotification(`获取任务列表失败: ${error.message}`, 'error');
            // 如果API调用失败，显示空列表
            updateTaskList([]);
        });
}

// 更新任务列表
function updateTaskList(tasks) {
    const taskListEl = document.getElementById('task-list');
    if (!taskListEl) return;
    
    taskListEl.innerHTML = '';
    
    if (!tasks || tasks.length === 0) {
        // 显示空状态
        const emptyRow = document.createElement('tr');
        emptyRow.innerHTML = `
            <td colspan="6" class="text-center text-base-content/60 py-8">
                <svg class="w-12 h-12 mx-auto mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
                </svg>
                暂无任务
            </td>
        `;
        taskListEl.appendChild(emptyRow);
        return;
    }
    
    tasks.forEach(task => {
        const row = createTaskRow(task);
        taskListEl.appendChild(row);
    });
}

// 创建任务行
function createTaskRow(task) {
    const row = document.createElement('tr');
    row.className = 'task-item';
    
    const statusBadge = getStatusBadge(task.status);
    const modeBadge = getModeBadge(task.mode);
    
    row.innerHTML = `
        <td>
            <div class="font-bold">${task.name}</div>
        </td>
        <td>${modeBadge}</td>
        <td>${statusBadge}</td>
        <td>${task.startTime}</td>
        <td>${task.endTime || '-'}</td>
        <td>
            <div class="flex gap-2">
                <button class="btn btn-sm btn-error btn-outline" onclick="deleteTask('${task.name}')" title="删除任务">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                    删除
                </button>
                <button class="btn btn-sm btn-primary btn-outline" onclick="viewTaskDetail('${task.name}')">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                    </svg>
                    详情
                </button>
            </div>
        </td>
    `;
    
    return row;
}

// 获取状态徽章
function getStatusBadge(status) {
    const statusMap = {
        'running': '<span class="badge badge-neutral gap-2"><svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-8.293l-3-3a1 1 0 00-1.414 1.414L10.586 9H7a1 1 0 100 2h3.586l-1.293 1.293a1 1 0 101.414 1.414l3-3a1 1 0 000-1.414z" clip-rule="evenodd"></path></svg>运行中</span>',
        'done': '<span class="badge badge-neutral gap-2"><svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path></svg>已完成</span>',
        'exit': '<span class="badge badge-neutral gap-2"><svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"></path></svg>已退出</span>',
        'completed': '<span class="badge badge-neutral gap-2"><svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path></svg>已完成</span>',
        'interrupted': '<span class="badge badge-neutral gap-2"><svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>中断</span>'
    };
    return statusMap[status] || statusMap['interrupted'];
}

// 获取运行模式徽章
function getModeBadge(mode) {
    if (!mode) {
        return '<span class="badge badge-neutral">未知</span>';
    }
    return `<span class="badge badge-neutral">${mode}</span>`;
}

// 获取任务头像
function getTaskAvatar(name) {
    const initials = name.substring(0, 2).toUpperCase();
    const colors = ['bg-primary', 'bg-secondary', 'bg-accent', 'bg-info', 'bg-success'];
    const color = colors[Math.floor(Math.random() * colors.length)];
    
    return `<div class="${color} text-${color.replace('bg-', '')}-content flex items-center justify-center">
        <span class="text-xs font-bold">${initials}</span>
    </div>`;
}

// 查看任务详情
function viewTaskDetail(taskName) {
    window.location.href = `task-detail.html?task=${encodeURIComponent(taskName)}`;
}

// 删除任务
// 全局变量存储要删除的任务名
let taskToDelete = null;

// 删除任务
function deleteTask(taskName) {
    taskToDelete = taskName;
    
    // 更新模态框内容
    const modalMessage = document.getElementById('delete-modal-message');
    modalMessage.textContent = `确定要删除任务 "${taskName}" 吗？`;
    
    // 显示模态框
    const modal = document.getElementById('delete-modal');
    modal.showModal();
}

// 关闭删除模态框
function closeDeleteModal() {
    const modal = document.getElementById('delete-modal');
    modal.close();
    taskToDelete = null;
}

// 确认删除
function confirmDelete() {
    if (!taskToDelete) return;
    
    const taskName = taskToDelete;
    taskToDelete = null;
    
    // 关闭模态框
    closeDeleteModal();
    
    // 找到对应的删除按钮并显示加载状态
    const taskRows = document.querySelectorAll('.task-item');
    let deleteButton = null;
    let originalContent = null;
    
    taskRows.forEach(row => {
        const taskNameCell = row.querySelector('td:first-child .font-bold');
        if (taskNameCell && taskNameCell.textContent === taskName) {
            deleteButton = row.querySelector('button[onclick*="deleteTask"]');
        }
    });
    
    if (deleteButton) {
        originalContent = deleteButton.innerHTML;
        deleteButton.innerHTML = `
            <span class="loading loading-spinner loading-sm"></span>
            删除中...
        `;
        deleteButton.disabled = true;
    }
    
    // 发送删除请求
    fetch(`/task/?task=${encodeURIComponent(taskName)}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // 删除成功，显示成功消息并刷新任务列表
            showNotification('任务删除成功', 'success');
            loadTaskList(); // 重新加载任务列表
        } else {
            // 删除失败，显示错误消息
            showNotification(data.message || '删除任务失败', 'error');
            // 恢复按钮状态
            if (deleteButton && originalContent) {
                deleteButton.innerHTML = originalContent;
                deleteButton.disabled = false;
            }
        }
    })
    .catch(error => {
        console.error('删除任务失败:', error);
        showNotification('删除任务失败: ' + error.message, 'error');
        // 恢复按钮状态
        if (deleteButton && originalContent) {
            deleteButton.innerHTML = originalContent;
            deleteButton.disabled = false;
        }
    });
}

// 设置事件监听器
function setupEventListeners() {
    // 导航栏菜单切换
    const menuToggle = document.querySelector('.dropdown-toggle');
    if (menuToggle) {
        menuToggle.addEventListener('click', function() {
            const dropdown = this.closest('.dropdown');
            dropdown.classList.toggle('dropdown-open');
        });
    }
    
    
    // 响应式菜单
    setupResponsiveMenu();
}

// 设置响应式菜单
function setupResponsiveMenu() {
    const mobileMenuToggle = document.querySelector('.btn-ghost.lg\\:hidden');
    if (mobileMenuToggle) {
        mobileMenuToggle.addEventListener('click', function() {
            const dropdown = this.closest('.dropdown');
            dropdown.classList.toggle('dropdown-open');
        });
    }
}


// 显示通知
function showNotification(message, type = 'info') {
    // 创建toast容器（如果不存在）
    let toastContainer = document.getElementById('toast-container');
    if (!toastContainer) {
        toastContainer = document.createElement('div');
        toastContainer.id = 'toast-container';
        toastContainer.className = 'toast toast-top toast-end z-50';
        document.body.appendChild(toastContainer);
    }
    
    const notification = document.createElement('div');
    notification.className = `alert alert-${type} max-w-sm`;
    
    // 根据类型选择图标
    let icon = '';
    switch(type) {
        case 'success':
            icon = '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>';
            break;
        case 'error':
            icon = '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>';
            break;
        case 'warning':
            icon = '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path></svg>';
            break;
        default:
            icon = '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>';
    }
    
    notification.innerHTML = `
        ${icon}
        <span>${message}</span>
        <button class="btn btn-sm btn-circle btn-ghost" onclick="this.parentElement.remove()">✕</button>
    `;
    
    toastContainer.appendChild(notification);
    
    // 自动移除通知
    setTimeout(() => {
        if (notification.parentElement) {
            notification.remove();
        }
    }, 5000);
}

// 显示加载状态
function showLoading(element, text = '加载中...') {
    const originalContent = element.innerHTML;
    element.innerHTML = `<span class="loading loading-spinner loading-sm"></span>${text}`;
    element.disabled = true;
    
    return function hideLoading() {
        element.innerHTML = originalContent;
        element.disabled = false;
    };
}

// 格式化时间
function formatTime(timestamp) {
    const date = new Date(timestamp);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    });
}

// 格式化文件大小
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 复制到剪贴板
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showNotification('已复制到剪贴板', 'success');
    }).catch(() => {
        showNotification('复制失败', 'error');
    });
}

// 防抖函数
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// 节流函数
function throttle(func, limit) {
    let inThrottle;
    return function() {
        const args = arguments;
        const context = this;
        if (!inThrottle) {
            func.apply(context, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}


// 错误处理
function handleError(error, context = '') {
    console.error(`Error in ${context}:`, error);
    showNotification(`操作失败: ${error.message}`, 'error');
}

// ==================== 命令页面相关功能 ====================

// 页面加载时获取帮助信息
function initializeCommandsPage() {
    loadHelpInfo();
}

// 加载帮助信息
function loadHelpInfo() {
    fetch('/execHelp')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.success && data.data) {
                // 更新模式选择器 - 修复：使用 flows 而不是 Modes
                updateModeSelect(data.data.flows);
                
                // 更新帮助信息 - 修复：使用 help 而不是 Help
                updateHelpContent(data.data.help);
            }
        })
        .catch(error => {
            console.error('获取帮助信息失败:', error);
        });
}

// 更新模式选择器
function updateModeSelect(modes) {
    const modeSelect = document.getElementById('mode-select');
    if (!modeSelect) return;
    
    modeSelect.innerHTML = '';
    
    // 确保 modes 是数组
    if (!Array.isArray(modes)) {
        modes = [];
    }
    
    if (modes.length === 0) {
        const option = document.createElement('option');
        option.value = '';
        option.textContent = '暂无可用工作流';
        modeSelect.appendChild(option);
        return;
    }
    
    modes.forEach(mode => {
        const option = document.createElement('option');
        option.value = mode;
        option.textContent = mode;
        // 默认选择第一个模式
        if (modes.indexOf(mode) === 0) {
            option.selected = true;
        }
        modeSelect.appendChild(option);
    });
}

// 更新帮助信息内容
function updateHelpContent(helpText) {
    const helpContent = document.getElementById('help-content');
    if (!helpContent) return;
    
    helpContent.textContent = helpText;
}

// 生成命令
function generateCommand() {
    const target = document.getElementById('target-input').value.trim();
    const mode = document.getElementById('mode-select').value;
    const taskName = document.getElementById('task-name-input').value.trim();
    const params = document.getElementById('params-input').value.trim();
    
    if (!target) {
        showNotification('请输入目标', 'warning');
        return;
    }
    
    let command = `goauto scan --target ${target} --flow ${mode}`;
    
    // 如果提供了任务名称，添加到命令中
    if (taskName) {
        command += ` --task-name ${taskName}`;
    }
    
    // 添加额外参数
    if (params) {
        command += ` ${params}`;
    }
    
    const generatedCommandEl = document.getElementById('generated-command');
    if (generatedCommandEl) {
        generatedCommandEl.value = command;
    }
}

// 执行命令
function executeCommand(event) {
    const command = document.getElementById('generated-command').value;
    if (!command.trim()) {
        showNotification('请先生成命令', 'warning');
        return;
    }

    // 显示执行中的状态
    const executeBtn = event ? event.target : document.querySelector('button[onclick*="executeCommand"]');
    const originalText = executeBtn.innerHTML;
    executeBtn.innerHTML = '<span class="loading loading-spinner loading-sm"></span>执行中...';
    executeBtn.disabled = true;

    // 调用后端API执行命令
    fetch(`/exec?command=${encodeURIComponent(command)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                showNotification('命令执行成功！', 'success');
                console.log('命令执行成功:', data);
            } else {
                throw new Error(data.message || '命令执行失败');
            }
        })
        .catch(error => {
            console.error('命令执行失败:', error);
            showNotification(`命令执行失败: ${error.message}`, 'error');
        })
        .finally(() => {
            // 恢复按钮状态
            executeBtn.innerHTML = originalText;
            executeBtn.disabled = false;
        });
}

// 上传目标
function uploadTargets() {
    const targetList = document.getElementById('target-list');
    if (!targetList) {
        showNotification('找不到目标列表输入框', 'error');
        return;
    }
    
    const targets = targetList.value.trim();
    if (!targets) {
        showNotification('请输入目标列表', 'warning');
        return;
    }

    // 显示上传中的状态
    const uploadBtn = document.querySelector('button[onclick*="uploadTargets"]');
    const originalText = uploadBtn.innerHTML;
    uploadBtn.innerHTML = '<span class="loading loading-spinner loading-sm"></span>上传中...';
    uploadBtn.disabled = true;

    // 创建FormData对象
    const formData = new FormData();
    formData.append('targets', targets);

    // 调用后端API上传目标
    fetch('/upload/targets', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            // 显示成功消息
            showNotification('目标上传成功！', 'success');
            
            // 更新文件路径显示
            const filePathEl = document.getElementById('file-path');
            const uploadResultEl = document.getElementById('upload-result');
            
            if (filePathEl) {
                filePathEl.textContent = data.data.filePath || '未知路径';
            }
            if (uploadResultEl) {
                uploadResultEl.style.display = 'block';
            }
            
            console.log('目标上传成功:', data);
        } else {
            throw new Error(data.message || '目标上传失败');
        }
    })
    .catch(error => {
        console.error('目标上传失败:', error);
        showNotification(`目标上传失败: ${error.message}`, 'error');
    })
    .finally(() => {
        // 恢复按钮状态
        uploadBtn.innerHTML = originalText;
        uploadBtn.disabled = false;
    });
}

// ==================== 页面初始化逻辑 ====================

// 检测当前页面并初始化相应功能
document.addEventListener('DOMContentLoaded', function() {
    const currentPage = window.location.pathname;
    
    if (currentPage.includes('commands.html')) {
        // 命令页面初始化
        initializeCommandsPage();
    } else {
        // 主页初始化
        initializeApp();
    }
});

// 导出函数供HTML使用
window.viewTaskDetail = viewTaskDetail;
window.showNotification = showNotification;
window.copyToClipboard = copyToClipboard;
window.generateCommand = generateCommand;
window.executeCommand = executeCommand;
window.uploadTargets = uploadTargets;
