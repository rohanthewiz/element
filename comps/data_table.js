document.addEventListener('DOMContentLoaded', function() {
    // Sample data (in real app, this would come from an API)
    const sampleData = [
        { id: 1, name: 'John Doe', email: 'john.doe@example.com', status: 'active', country: 'USA', created: '2025-01-15' },
        { id: 2, name: 'Jane Smith', email: 'jane.smith@example.com', status: 'inactive', country: 'Canada', created: '2025-01-18' },
        { id: 3, name: 'Bob Johnson', email: 'bob.johnson@example.com', status: 'active', country: 'UK', created: '2025-01-20' },
        { id: 4, name: 'Alice Williams', email: 'alice.williams@example.com', status: 'pending', country: 'Germany', created: '2025-01-22' },
        { id: 5, name: 'Charlie Brown', email: 'charlie.brown@example.com', status: 'active', country: 'France', created: '2025-01-25' },
        { id: 6, name: 'Diana Miller', email: 'diana.miller@example.com', status: 'inactive', country: 'USA', created: '2025-01-27' },
        { id: 7, name: 'Edward Davis', email: 'edward.davis@example.com', status: 'active', country: 'Canada', created: '2025-01-28' },
        { id: 8, name: 'Fiona Clark', email: 'fiona.clark@example.com', status: 'pending', country: 'UK', created: '2025-01-30' },
        { id: 9, name: 'George Wilson', email: 'george.wilson@example.com', status: 'active', country: 'Germany', created: '2025-02-01' },
        { id: 10, name: 'Hannah Moore', email: 'hannah.moore@example.com', status: 'inactive', country: 'France', created: '2025-02-03' },
        { id: 11, name: 'Ian Taylor', email: 'ian.taylor@example.com', status: 'active', country: 'Japan', created: '2025-02-05' },
        { id: 12, name: 'Julia Adams', email: 'julia.adams@example.com', status: 'pending', country: 'Australia', created: '2025-02-07' },
        { id: 13, name: 'Kevin White', email: 'kevin.white@example.com', status: 'active', country: 'Brazil', created: '2025-02-09' },
        { id: 14, name: 'Laura Harris', email: 'laura.harris@example.com', status: 'inactive', country: 'India', created: '2025-02-11' },
        { id: 15, name: 'Mike Robinson', email: 'mike.robinson@example.com', status: 'active', country: 'China', created: '2025-02-13' },
        { id: 16, name: 'Nancy Garcia', email: 'nancy.garcia@example.com', status: 'pending', country: 'USA', created: '2025-02-15' },
        { id: 17, name: 'Oscar Lee', email: 'oscar.lee@example.com', status: 'active', country: 'Canada', created: '2025-02-17' },
        { id: 18, name: 'Patricia King', email: 'patricia.king@example.com', status: 'inactive', country: 'UK', created: '2025-02-19' },
        { id: 19, name: 'Quincy Martinez', email: 'quincy.martinez@example.com', status: 'active', country: 'Germany', created: '2025-02-21' },
        { id: 20, name: 'Rachel Scott', email: 'rachel.scott@example.com', status: 'pending', country: 'France', created: '2025-02-23' }
    ];

    // State variables
    let data = [...sampleData];
    let currentPage = 1;
    let itemsPerPage = 10;
    let sortColumn = 'id';
    let sortDirection = 'asc';
    let searchTerm = '';
    let filters = {};
    let selectedRows = new Set();
    let visibleColumns = {
        id: true,
        name: true,
        email: true,
        status: true,
        country: true,
        created: true,
        actions: true
    };

    // DOM Elements
    const dataTable = document.getElementById('$DATA_TABLE_ID');
    const tableBody = dataTable.querySelector('tbody');
    const searchInput = document.getElementById('searchInput');
    const perPageSelector = document.getElementById('perPageSelector');
    const paginationControls = document.getElementById('paginationControls');
    const showingStart = document.getElementById('showingStart');
    const showingEnd = document.getElementById('showingEnd');
    const totalItems = document.getElementById('totalItems');
    const exportBtn = document.getElementById('exportBtn');
    const selectAll = document.getElementById('selectAll');
    const bulkActions = document.getElementById('bulkActions');
    const selectedCount = document.getElementById('selectedCount');
    const bulkDeleteBtn = document.getElementById('bulkDeleteBtn');
    const bulkStatusBtn = document.getElementById('bulkStatusBtn');
    const bulkCancelBtn = document.getElementById('bulkCancelBtn');
    const columnToggleBtn = document.getElementById('columnToggleBtn');
    const columnToggleDropdown = document.getElementById('columnToggleDropdown');
    const loadingOverlay = document.querySelector('.loading-overlay');
    const userModal = document.getElementById('userModal');
    const modalTitle = document.getElementById('modalTitle');
    const modalClose = document.getElementById('modalClose');
    const userForm = document.getElementById('userForm');
    const userId = document.getElementById('userId');
    const userName = document.getElementById('userName');
    const userEmail = document.getElementById('userEmail');
    const userStatus = document.getElementById('userStatus');
    const userCountry = document.getElementById('userCountry');
    const cancelBtn = document.getElementById('cancelBtn');
    const saveBtn = document.getElementById('saveBtn');
    const viewModal = document.getElementById('viewModal');
    const viewModalClose = document.getElementById('viewModalClose');
    const viewModalBody = document.getElementById('viewModalBody');
    const closeViewBtn = document.getElementById('closeViewBtn');
    const confirmModal = document.getElementById('confirmModal');
    const confirmModalClose = document.getElementById('confirmModalClose');
    const confirmMessage = document.getElementById('confirmMessage');
    const cancelConfirmBtn = document.getElementById('cancelConfirmBtn');
    const confirmBtn = document.getElementById('confirmBtn');
    const toast = document.getElementById('toast');
    const nameError = document.getElementById('nameError');
    const emailError = document.getElementById('emailError');
    const lastUpdated = document.getElementById('lastUpdated');

    // Initialize
    initColumnToggleDropdown();
    updateLastUpdated();
    renderData();

    // Event Listeners
    searchInput.addEventListener('input', handleSearch);
    perPageSelector.addEventListener('change', handlePerPageChange);
    exportBtn.addEventListener('click', exportToCsv);
    selectAll.addEventListener('change', handleSelectAll);
    bulkDeleteBtn.addEventListener('click', handleBulkDelete);
    bulkStatusBtn.addEventListener('click', handleBulkStatus);
    bulkCancelBtn.addEventListener('click', handleBulkCancel);
    columnToggleBtn.addEventListener('click', toggleColumnDropdown);

    // Modal event listeners
    modalClose.addEventListener('click', closeUserModal);
    cancelBtn.addEventListener('click', closeUserModal);
    saveBtn.addEventListener('click', saveUser);
    viewModalClose.addEventListener('click', closeViewModal);
    closeViewBtn.addEventListener('click', closeViewModal);
    confirmModalClose.addEventListener('click', closeConfirmModal);
    cancelConfirmBtn.addEventListener('click', closeConfirmModal);

    // Table header click for sorting
    document.querySelectorAll('th.sortable').forEach(th => {
        th.addEventListener('click', () => {
            const column = th.dataset.column;
            if (column) {
                handleSort(column);
            }
        });
    });

    // Click outside to close dropdowns
    document.addEventListener('click', function(event) {
        if (!columnToggleBtn.contains(event.target) && !columnToggleDropdown.contains(event.target)) {
            columnToggleDropdown.style.display = 'none';
        }

        // Close any open column filters
        const filterDropdowns = document.querySelectorAll('.column-filter');
        filterDropdowns.forEach(dropdown => {
            if (dropdown.style.display === 'block' && !dropdown.contains(event.target)) {
                dropdown.style.display = 'none';
            }
        });
    });

    // Functions

    // Initialize column toggle dropdown
    function initColumnToggleDropdown() {
        columnToggleDropdown.innerHTML = '';

        Object.keys(visibleColumns).forEach(column => {
            if (column !== 'actions') { // Always keep actions column
                const label = document.createElement('label');
                label.className = 'filter-item';

                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.checked = visibleColumns[column];
                checkbox.addEventListener('change', () => {
                    visibleColumns[column] = checkbox.checked;
                    renderData();
                });

                label.appendChild(checkbox);
                label.appendChild(document.createTextNode(' ' + formatColumnName(column)));
                columnToggleDropdown.appendChild(label);
            }
        });
    }

    // Format column name for display
    function formatColumnName(column) {
        return column.charAt(0).toUpperCase() + column.slice(1);
    }

    // Toggle column dropdown
    function toggleColumnDropdown() {
        columnToggleDropdown.style.display = columnToggleDropdown.style.display === 'block' ? 'none' : 'block';
    }

    // Filter, sort and paginate data
    function getProcessedData() {
        // Apply search
        let result = data.filter(item => {
            if (!searchTerm) return true;

            // Search in all text fields
            return Object.values(item).some(value =>
                String(value).toLowerCase().includes(searchTerm.toLowerCase())
            );
        });

        // Apply column filters
        Object.keys(filters).forEach(column => {
            const filterValues = filters[column];

            if (filterValues && filterValues.length > 0) {
                result = result.filter(item =>
                    filterValues.includes(String(item[column]))
                );
            }
        });

        // Apply sorting
        result.sort((a, b) => {
            let aValue = a[sortColumn];
            let bValue = b[sortColumn];

            // Handle string comparison
            if (typeof aValue === 'string') {
                aValue = aValue.toLowerCase();
                bValue = bValue.toLowerCase();
            }

            if (aValue < bValue) {
                return sortDirection === 'asc' ? -1 : 1;
            }
            if (aValue > bValue) {
                return sortDirection === 'asc' ? 1 : -1;
            }
            return 0;
        });

        return result;
    }

    // Render data to table
    function renderData() {
        showLoading();

        // Small delay to show loading effect
        setTimeout(() => {
            const processedData = getProcessedData();
            totalItems.textContent = processedData.length;

            // Calculate pagination
            const startIndex = (currentPage - 1) * itemsPerPage;
            const endIndex = Math.min(startIndex + itemsPerPage, processedData.length);
            const paginatedData = processedData.slice(startIndex, endIndex);

            // Update showing info
            showingStart.textContent = processedData.length ? startIndex + 1 : 0;
            showingEnd.textContent = endIndex;

            // Clear table body
            tableBody.innerHTML = '';

            // Render rows
            if (paginatedData.length === 0) {
                // No data message
                const noDataRow = document.createElement('tr');
                const noDataCell = document.createElement('td');
                noDataCell.colSpan = 8;
                noDataCell.textContent = 'No data found';
                noDataCell.style.textAlign = 'center';
                noDataCell.style.padding = '20px';
                noDataRow.appendChild(noDataCell);
                tableBody.appendChild(noDataRow);
            } else {
                paginatedData.forEach(item => {
                    const row = document.createElement('tr');
                    if (selectedRows.has(item.id)) {
                        row.classList.add('row-selected');
                    }

                    // Checkbox cell
                    const checkboxCell = document.createElement('td');
                    const checkbox = document.createElement('input');
                    checkbox.type = 'checkbox';
                    checkbox.className = 'row-checkbox';
                    checkbox.checked = selectedRows.has(item.id);
                    checkbox.addEventListener('change', () => {
                        if (checkbox.checked) {
                            selectedRows.add(item.id);
                            row.classList.add('row-selected');
                        } else {
                            selectedRows.delete(item.id);
                            row.classList.remove('row-selected');
                        }
                        updateBulkActions();
                    });
                    checkboxCell.appendChild(checkbox);
                    row.appendChild(checkboxCell);

                    // Data cells
                    Object.keys(visibleColumns).forEach(column => {
                        if (column !== 'actions' && visibleColumns[column]) {
                            const cell = document.createElement('td');

                            if (column === 'status') {
                                // Status badge
                                const badge = document.createElement('span');
                                badge.className = `status-badge status-${item[column]}`;
                                badge.textContent = item[column].charAt(0).toUpperCase() + item[column].slice(1);
                                cell.appendChild(badge);
                            } else {
                                // Regular cell with text
                                cell.textContent = item[column];
                            }

                            row.appendChild(cell);
                        }
                    });

                    // Actions cell
                    if (visibleColumns.actions) {
                        const actionsCell = document.createElement('td');
                        actionsCell.className = 'action-cell';

                        // View button
                        const viewBtn = document.createElement('button');
                        viewBtn.className = 'action-btn view';
                        viewBtn.textContent = 'View';
                        viewBtn.addEventListener('click', () => openViewModal(item));
                        actionsCell.appendChild(viewBtn);

                        // Edit button
                        const editBtn = document.createElement('button');
                        editBtn.className = 'action-btn edit';
                        editBtn.textContent = 'Edit';
                        editBtn.addEventListener('click', () => openEditModal(item));
                        actionsCell.appendChild(editBtn);

                        // Delete button
                        const deleteBtn = document.createElement('button');
                        deleteBtn.className = 'action-btn delete';
                        deleteBtn.textContent = 'Delete';
                        deleteBtn.addEventListener('click', () => openDeleteConfirmation(item));
                        actionsCell.appendChild(deleteBtn);

                        row.appendChild(actionsCell);
                    }

                    tableBody.appendChild(row);
                });
            }

            // Update pagination controls
            renderPagination(processedData.length);

            // Update "select all" checkbox state
            updateSelectAllCheckbox();

            hideLoading();
        }, 300);
    }

    // Render pagination controls
    function renderPagination(totalCount) {
        paginationControls.innerHTML = '';

        if (totalCount === 0) {
            return;
        }

        const totalPages = Math.ceil(totalCount / itemsPerPage);

        // Previous button
        const prevBtn = document.createElement('button');
        prevBtn.className = 'pagination-btn';
        prevBtn.textContent = '←';
        prevBtn.disabled = currentPage === 1;
        prevBtn.addEventListener('click', () => {
            if (currentPage > 1) {
                currentPage--;
                renderData();
            }
        });
        paginationControls.appendChild(prevBtn);

        // Page buttons
        let startPage = Math.max(1, currentPage - 2);
        let endPage = Math.min(totalPages, startPage + 4);

        if (endPage - startPage < 4) {
            startPage = Math.max(1, endPage - 4);
        }

        for (let i = startPage; i <= endPage; i++) {
            const pageBtn = document.createElement('button');
            pageBtn.className = 'pagination-btn';
            if (i === currentPage) {
                pageBtn.classList.add('active');
            }
            pageBtn.textContent = i;
            pageBtn.addEventListener('click', () => {
                currentPage = i;
                renderData();
            });
            paginationControls.appendChild(pageBtn);
        }

        // Next button
        const nextBtn = document.createElement('button');
        nextBtn.className = 'pagination-btn';
        nextBtn.textContent = '→';
        nextBtn.disabled = currentPage === totalPages;
        nextBtn.addEventListener('click', () => {
            if (currentPage < totalPages) {
                currentPage++;
                renderData();
            }
        });
        paginationControls.appendChild(nextBtn);
    }

    // Handle sort by column
    function handleSort(column) {
        if (sortColumn === column) {
            // Toggle direction if same column
            sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            // New column, default to ascending
            sortColumn = column;
            sortDirection = 'asc';
        }

        // Update UI
        document.querySelectorAll('th.sortable').forEach(th => {
            th.classList.remove('asc', 'desc');
        });

        const th = document.querySelector(`th[data-column="${column}"]`);
        th.classList.add(sortDirection);

        // Reset to first page and render
        currentPage = 1;
        renderData();
    }

    // Handle search
    function handleSearch() {
        searchTerm = searchInput.value;
        currentPage = 1;
        renderData();
    }

    // Handle items per page change
    function handlePerPageChange() {
        itemsPerPage = parseInt(perPageSelector.value);
        currentPage = 1;
        renderData();
    }

    // Export to CSV
    function exportToCsv() {
        const processedData = getProcessedData();

        if (processedData.length === 0) {
            showToast('No data to export', 'warning');
            return;
        }

        // Get visible columns
        const columns = Object.keys(visibleColumns)
            .filter(column => visibleColumns[column] && column !== 'actions');

        // Create headers row
        let csvContent = columns.map(formatColumnName).join(',') + '\n';

        // Create data rows
        processedData.forEach(item => {
            const row = columns.map(column => {
                let value = item[column];

                // Escape quotes and handle commas
                if (typeof value === 'string') {
                    value = value.replace(/"/g, '""');
                    if (value.includes(',')) {
                        value = `"${value}"`;
                    }
                }

                return value;
            }).join(',');

            csvContent += row + '\n';
        });

        // Create download link
        const encodedUri = 'data:text/csv;charset=utf-8,' + encodeURIComponent(csvContent);
        const link = document.createElement('a');
        link.setAttribute('href', encodedUri);
        link.setAttribute('download', 'customer_data_export.csv');
        document.body.appendChild(link);

        // Trigger download and clean up
        link.click();
        document.body.removeChild(link);

        showToast('Export successful!', 'success');
    }

    // Select all rows
    function handleSelectAll() {
        const processedData = getProcessedData();
        const checkboxes = document.querySelectorAll('.row-checkbox');

        if (selectAll.checked) {
            // Select all visible rows
            checkboxes.forEach((checkbox, index) => {
                const startIndex = (currentPage - 1) * itemsPerPage;
                const item = processedData[startIndex + index];
                if (item) {
                    selectedRows.add(item.id);
                    checkbox.checked = true;
                    checkbox.closest('tr').classList.add('row-selected');
                }
            });
        } else {
            // Deselect all rows
            checkboxes.forEach(checkbox => {
                checkbox.checked = false;
                checkbox.closest('tr').classList.remove('row-selected');
            });
            selectedRows.clear();
        }

        updateBulkActions();
    }

    // Update select all checkbox state
    function updateSelectAllCheckbox() {
        const checkboxes = document.querySelectorAll('.row-checkbox');
        const checkedCount = Array.from(checkboxes).filter(cb => cb.checked).length;

        if (checkboxes.length > 0 && checkedCount === checkboxes.length) {
            selectAll.checked = true;
            selectAll.indeterminate = false;
        } else if (checkedCount > 0) {
            selectAll.checked = false;
            selectAll.indeterminate = true;
        } else {
            selectAll.checked = false;
            selectAll.indeterminate = false;
        }
    }

    // Update bulk actions display
    function updateBulkActions() {
        selectedCount.textContent = selectedRows.size;

        if (selectedRows.size > 0) {
            bulkActions.style.display = 'flex';
        } else {
            bulkActions.style.display = 'none';
        }
    }

    // Handle bulk delete
    function handleBulkDelete() {
        if (selectedRows.size === 0) return;

        confirmMessage.textContent = `Are you sure you want to delete ${selectedRows.size} selected items?`;

        openConfirmModal(() => {
            // Perform delete
            data = data.filter(item => !selectedRows.has(item.id));
            selectedRows.clear();
            updateBulkActions();
            renderData();
            updateLastUpdated();

            showToast(`${selectedRows.size} items deleted successfully`, 'success');
        });
    }

    // Handle bulk status change
    function handleBulkStatus() {
        if (selectedRows.size === 0) return;

        const statuses = ['active', 'inactive', 'pending'];
        let modalContent = `
                    <p>Change status for ${selectedRows.size} selected items to:</p>
                    <div class="form-group">
                        <select class="form-select" id="bulkStatusSelect">
                            ${statuses.map(status => `<option value="${status}">${status.charAt(0).toUpperCase() + status.slice(1)}</option>`).join('')}
                        </select>
                    </div>
                `;

        confirmMessage.innerHTML = modalContent;

        openConfirmModal(() => {
            const newStatus = document.getElementById('bulkStatusSelect').value;

            // Update statuses
            data.forEach(item => {
                if (selectedRows.has(item.id)) {
                    item.status = newStatus;
                }
            });

            renderData();
            updateLastUpdated();

            showToast(`Status updated for ${selectedRows.size} items`, 'success');
        });
    }

    // Handle bulk cancel
    function handleBulkCancel() {
        selectedRows.clear();
        updateBulkActions();
        renderData();
    }

    // Show loading overlay
    function showLoading() {
        loadingOverlay.style.display = 'flex';
    }

    // Hide loading overlay
    function hideLoading() {
        loadingOverlay.style.display = 'none';
    }

    // Open user modal for editing
    function openEditModal(item) {
        modalTitle.textContent = 'Edit User';
        userId.value = item.id;
        userName.value = item.name;
        userEmail.value = item.email;
        userStatus.value = item.status;
        userCountry.value = item.country;

        clearErrors();
        userModal.style.display = 'flex';
    }

    // Open user modal for adding
    function openAddModal() {
        modalTitle.textContent = 'Add New User';
        userForm.reset();
        userId.value = '';

        clearErrors();
        userModal.style.display = 'flex';
    }

    // Close user modal
    function closeUserModal() {
        userModal.style.display = 'none';
    }

    // Save user (add/edit)
    function saveUser() {
        // Validate form
        if (!validateForm()) {
            return;
        }

        const isEditing = userId.value !== '';

        if (isEditing) {
            // Edit existing user
            const id = parseInt(userId.value);
            const index = data.findIndex(item => item.id === id);

            if (index !== -1) {
                data[index] = {
                    ...data[index],
                    name: userName.value,
                    email: userEmail.value,
                    status: userStatus.value,
                    country: userCountry.value
                };

                showToast('User updated successfully', 'success');
            }
        } else {
            // Add new user
            const newId = Math.max(...data.map(item => item.id), 0) + 1;

            data.push({
                id: newId,
                name: userName.value,
                email: userEmail.value,
                status: userStatus.value,
                country: userCountry.value,
                created: new Date().toISOString().split('T')[0]
            });

            showToast('User added successfully', 'success');
        }

        closeUserModal();
        updateLastUpdated();
        renderData();
    }

    // Validate form
    function validateForm() {
        let isValid = true;
        clearErrors();

        if (!userName.value.trim()) {
            nameError.textContent = 'Name is required';
            isValid = false;
        }

        if (!userEmail.value.trim()) {
            emailError.textContent = 'Email is required';
            isValid = false;
        } else if (!isValidEmail(userEmail.value)) {
            emailError.textContent = 'Please enter a valid email address';
            isValid = false;
        }

        return isValid;
    }

    // Clear form errors
    function clearErrors() {
        nameError.textContent = '';
        emailError.textContent = '';
    }

    // Validate email format
    function isValidEmail(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    }

    // Open view modal
    function openViewModal(item) {
        let html = '<div style="line-height: 1.6;">';

        Object.keys(item).forEach(key => {
            if (key === 'status') {
                html += `<p><strong>${formatColumnName(key)}:</strong> <span class="status-badge status-${item[key]}">${item[key].charAt(0).toUpperCase() + item[key].slice(1)}</span></p>`;
            } else {
                html += `<p><strong>${formatColumnName(key)}:</strong> ${item[key]}</p>`;
            }
        });

        html += '</div>';

        viewModalBody.innerHTML = html;
        viewModal.style.display = 'flex';
    }

    // Close view modal
    function closeViewModal() {
        viewModal.style.display = 'none';
    }

    // Open delete confirmation
    function openDeleteConfirmation(item) {
        confirmMessage.textContent = `Are you sure you want to delete "${item.name}"?`;

        openConfirmModal(() => {
            // Remove item
            data = data.filter(dataItem => dataItem.id !== item.id);

            // Also remove from selected if present
            if (selectedRows.has(item.id)) {
                selectedRows.delete(item.id);
                updateBulkActions();
            }

            renderData();
            updateLastUpdated();

            showToast('Item deleted successfully', 'success');
        });
    }

    // Open confirmation modal
    function openConfirmModal(onConfirm) {
        confirmModal.style.display = 'flex';

        // Set confirm button action
        confirmBtn.onclick = () => {
            onConfirm();
            closeConfirmModal();
        };
    }

    // Close confirmation modal
    function closeConfirmModal() {
        confirmModal.style.display = 'none';
    }

    // Show toast notification
    function showToast(message, type = 'success') {
        toast.textContent = message;
        toast.className = 'toast';
        toast.classList.add(type);
        toast.style.display = 'block';

        setTimeout(() => {
            toast.style.display = 'none';
        }, 3000);
    }

    // Update last updated timestamp
    function updateLastUpdated() {
        const now = new Date();
        lastUpdated.textContent = now.toLocaleString();
    }

    // Add Row Button (at the bottom of the page)
    const addRowButton = document.createElement('button');
    addRowButton.className = 'export-btn';
    addRowButton.style.marginTop = '20px';
    addRowButton.style.marginLeft = '20px';
    addRowButton.textContent = '+ Add New User';
    addRowButton.addEventListener('click', openAddModal);
    document.querySelector('.data-table-container').after(addRowButton);
});
