/* ============================================
   FORUM — VANILLA JS
   ============================================ */

// ============================================
// TOAST NOTIFICATION SYSTEM
// ============================================
const Toast = {
    container: null,
    queue: [],
    maxVisible: 3,

    init() {
        this.container = document.getElementById('toastContainer');
    },

    show(message, type) {
        if (!this.container) this.init();
        const visible = this.container.querySelectorAll('.toast:not(.removing)');
        if (visible.length >= this.maxVisible) {
            this.queue.push({ message, type });
            return;
        }

        const icons = {
            success: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>',
            error: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>',
            info: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>'
        };

        const toast = document.createElement('div');
        toast.className = 'toast ' + type;
        toast.innerHTML =
            '<span class="toast-icon">' + (icons[type] || icons.info) + '</span>' +
            '<span class="toast-message">' + message + '</span>' +
            '<button class="toast-close" onclick="Toast.remove(this.parentElement)">&times;</button>';

        this.container.appendChild(toast);

        setTimeout(() => this.remove(toast), 3000);
    },

    remove(el) {
        if (!el || el.classList.contains('removing')) return;
        el.classList.add('removing');
        setTimeout(() => {
            el.remove();
            if (this.queue.length > 0) {
                const next = this.queue.shift();
                this.show(next.message, next.type);
            }
        }, 300);
    },

    success(msg) { this.show(msg, 'success'); },
    error(msg) { this.show(msg, 'error'); },
    info(msg) { this.show(msg, 'info'); }
};

// ============================================
// THEME TOGGLE
// ============================================
(function() {
    var toggle = document.getElementById('themeToggle');
    if (!toggle) return;
    toggle.addEventListener('click', function() {
        var current = document.documentElement.getAttribute('data-theme') || 'dark';
        var next = current === 'dark' ? 'light' : 'dark';
        document.documentElement.setAttribute('data-theme', next);
        localStorage.setItem('theme', next);
    });
})();

// ============================================
// HAMBURGER MENU
// ============================================
(function() {
    const hamburger = document.getElementById('hamburger');
    const nav = document.getElementById('nav');
    if (hamburger && nav) {
        hamburger.addEventListener('click', function() {
            hamburger.classList.toggle('active');
            nav.classList.toggle('active');
        });
    }
})();

// ============================================
// USER DROPDOWN
// ============================================
(function() {
    const navUser = document.getElementById('navUser');
    const dropdown = document.getElementById('userDropdown');
    if (navUser && dropdown) {
        navUser.addEventListener('click', function(e) {
            e.stopPropagation();
            dropdown.classList.toggle('active');
        });
        document.addEventListener('click', function() {
            dropdown.classList.remove('active');
        });
    }
})();

// ============================================
// LIVE SEARCH (debounced)
// ============================================
(function() {
    const input = document.getElementById('searchInput');
    const dropdown = document.getElementById('searchDropdown');
    if (!input || !dropdown) return;

    let timer = null;

    input.addEventListener('input', function() {
        const q = this.value.trim();
        clearTimeout(timer);

        if (q.length < 2) {
            dropdown.classList.remove('active');
            dropdown.innerHTML = '';
            return;
        }

        timer = setTimeout(function() {
            fetch('/api/search?q=' + encodeURIComponent(q), {
                headers: { 'Accept': 'application/json' }
            })
            .then(function(r) { return r.json(); })
            .then(function(results) {
                if (!results || results.length === 0) {
                    dropdown.innerHTML = '<div style="padding:16px;color:var(--text-muted);text-align:center;font-size:0.9rem">Aucun résultat</div>';
                    dropdown.classList.add('active');
                    return;
                }

                let html = '';
                results.forEach(function(r) {
                    html += '<a href="/post/' + r.id + '" class="search-result">' +
                        '<div class="search-result-title">' + escapeHtml(r.title) + '</div>' +
                        '<div class="search-result-meta">par ' + escapeHtml(r.author) + '</div>' +
                        '</a>';
                });
                html += '<a href="/search?q=' + encodeURIComponent(q) + '" class="search-all-link">Voir tous les résultats</a>';
                dropdown.innerHTML = html;
                dropdown.classList.add('active');
            })
            .catch(function() {
                dropdown.classList.remove('active');
            });
        }, 300);
    });

    input.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            const q = this.value.trim();
            if (q) window.location = '/search?q=' + encodeURIComponent(q);
        }
    });

    document.addEventListener('click', function(e) {
        if (!e.target.closest('.nav-search')) {
            dropdown.classList.remove('active');
        }
    });

    input.addEventListener('focus', function() {
        if (dropdown.innerHTML.trim()) dropdown.classList.add('active');
    });
})();

// ============================================
// AJAX: LIKE POST
// ============================================
function handleLike(btn, type) {
    const postId = btn.dataset.postId;
    if (!postId) return;

    const formData = new FormData();
    formData.append('type', 'like');

    fetch('/post/react/' + postId, {
        method: 'POST',
        headers: { 'X-Requested-With': 'XMLHttpRequest' },
        body: formData
    })
    .then(function(r) {
        if (r.status === 401) {
            window.location = '/login';
            return null;
        }
        return r.json();
    })
    .then(function(data) {
        if (!data || !data.success) return;

        // Update all like buttons for this post on the page
        document.querySelectorAll('.like-action[data-post-id="' + postId + '"]').forEach(function(b) {
            if (data.userVote === 'like') {
                b.classList.add('active');
            } else {
                b.classList.remove('active');
            }
            const count = b.querySelector('.like-count');
            if (count) count.textContent = data.likes;
        });

        // Remove dislike active state
        document.querySelectorAll('.dislike-action[data-post-id="' + postId + '"]').forEach(function(b) {
            b.classList.remove('active');
            const count = b.querySelector('.dislike-count');
            if (count) count.textContent = data.dislikes;
        });

        // Heartbeat animation
        const svg = btn.querySelector('svg');
        if (svg && data.userVote === 'like') {
            svg.style.animation = 'heartbeat 0.4s ease';
            setTimeout(function() { svg.style.animation = ''; }, 400);
        }
    })
    .catch(function() {
        Toast.error('Erreur de connexion');
    });
}

// ============================================
// AJAX: REPOST
// ============================================
function handleRepost(btn) {
    const postId = btn.dataset.postId;
    if (!postId) return;

    fetch('/repost/' + postId, {
        method: 'POST',
        headers: { 'X-Requested-With': 'XMLHttpRequest' }
    })
    .then(function(r) {
        if (r.status === 401) {
            window.location = '/login';
            return null;
        }
        return r.json();
    })
    .then(function(data) {
        if (!data || !data.success) return;

        document.querySelectorAll('.repost-action[data-post-id="' + postId + '"]').forEach(function(b) {
            if (data.reposted) {
                b.classList.add('active');
                Toast.success('Reposté !');
            } else {
                b.classList.remove('active');
            }
            const count = b.querySelector('.repost-count');
            if (count) count.textContent = data.count;
        });

        // Rotate animation
        const svg = btn.querySelector('svg');
        if (svg && data.reposted) {
            svg.style.transition = 'transform 0.3s';
            svg.style.transform = 'rotate(360deg)';
            setTimeout(function() {
                svg.style.transition = '';
                svg.style.transform = '';
            }, 300);
        }
    })
    .catch(function() {
        Toast.error('Erreur de connexion');
    });
}

// ============================================
// AJAX: COMMENT REACTION
// ============================================
function handleCommentReaction(btn, type) {
    const commentId = btn.dataset.commentId;
    const postId = btn.dataset.postId;
    if (!commentId) return;

    const formData = new FormData();
    formData.append('type', type);
    formData.append('post_id', postId);

    fetch('/comment/react/' + commentId, {
        method: 'POST',
        headers: { 'X-Requested-With': 'XMLHttpRequest' },
        body: formData
    })
    .then(function(r) {
        if (r.status === 401) {
            window.location = '/login';
            return null;
        }
        return r.json();
    })
    .then(function(data) {
        if (!data || !data.success) return;

        const card = btn.closest('.comment-card') || btn.closest('.comment-actions');
        if (!card) return;

        const likeBtn = card.querySelector('.like-action');
        const dislikeBtn = card.querySelector('.dislike-action');

        if (likeBtn) {
            if (data.userVote === 'like') likeBtn.classList.add('active');
            else likeBtn.classList.remove('active');
            const lc = likeBtn.querySelector('.like-count');
            if (lc) lc.textContent = data.likes;
        }

        if (dislikeBtn) {
            if (data.userVote === 'dislike') dislikeBtn.classList.add('active');
            else dislikeBtn.classList.remove('active');
            const dc = dislikeBtn.querySelector('.dislike-count');
            if (dc) dc.textContent = data.dislikes;
        }
    })
    .catch(function() {
        Toast.error('Erreur de connexion');
    });
}

// ============================================
// SHARE DROPDOWN
// ============================================
function toggleShare(btn) {
    const dropdown = btn.parentElement.querySelector('.share-dropdown');
    if (!dropdown) return;

    // Close all other share dropdowns
    document.querySelectorAll('.share-dropdown.active').forEach(function(d) {
        if (d !== dropdown) d.classList.remove('active');
    });

    dropdown.classList.toggle('active');
}

function copyLink(path) {
    const url = window.location.origin + path;
    navigator.clipboard.writeText(url).then(function() {
        Toast.success('Lien copié !');
    }).catch(function() {
        // Fallback
        const input = document.createElement('input');
        input.value = url;
        document.body.appendChild(input);
        input.select();
        document.execCommand('copy');
        document.body.removeChild(input);
        Toast.success('Lien copié !');
    });
    closeAllShareDropdowns();
}

function shareTwitter(title, path) {
    const url = window.location.origin + path;
    window.open('https://twitter.com/intent/tweet?text=' + encodeURIComponent(title) + '&url=' + encodeURIComponent(url), '_blank', 'width=550,height=420');
    closeAllShareDropdowns();
}

function shareWhatsApp(title, path) {
    const url = window.location.origin + path;
    window.open('https://wa.me/?text=' + encodeURIComponent(title + ' ' + url), '_blank');
    closeAllShareDropdowns();
}

function closeAllShareDropdowns() {
    document.querySelectorAll('.share-dropdown.active').forEach(function(d) {
        d.classList.remove('active');
    });
}

// Close share dropdowns on outside click
document.addEventListener('click', function(e) {
    if (!e.target.closest('.share-action') && !e.target.closest('.share-dropdown')) {
        closeAllShareDropdowns();
    }
});

// ============================================
// COMMENT FORM — ENABLE/DISABLE SUBMIT
// ============================================
(function() {
    const textarea = document.getElementById('commentText');
    const submit = document.getElementById('commentSubmit');
    if (textarea && submit) {
        textarea.addEventListener('input', function() {
            submit.disabled = this.value.trim().length === 0;
        });
    }
})();

// ============================================
// REGISTER FORM VALIDATION
// ============================================
(function() {
    const form = document.getElementById('registerForm');
    if (!form) return;

    const username = document.getElementById('username');
    const email = document.getElementById('email');
    const password = document.getElementById('password');
    const confirm = document.getElementById('confirm_password');
    const strengthFill = document.getElementById('strengthFill');
    const strengthText = document.getElementById('strengthText');

    function showError(id, msg) {
        const el = document.getElementById(id);
        if (el) { el.textContent = msg; el.style.display = msg ? 'block' : 'none'; }
    }

    if (username) {
        username.addEventListener('input', function() {
            const v = this.value;
            if (v.length > 0 && v.length < 3) {
                showError('usernameError', 'Min. 3 caractères');
            } else if (/\s/.test(v)) {
                showError('usernameError', 'Pas d\'espaces');
            } else if (v.length > 20) {
                showError('usernameError', 'Max. 20 caractères');
            } else {
                showError('usernameError', '');
            }
        });
    }

    if (email) {
        email.addEventListener('input', function() {
            const v = this.value;
            if (v.length > 0 && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)) {
                showError('emailError', 'Email invalide');
            } else {
                showError('emailError', '');
            }
        });
    }

    if (password) {
        password.addEventListener('input', function() {
            const v = this.value;
            let strength = 0;
            if (v.length >= 8) strength++;
            if (/[A-Z]/.test(v) && /[a-z]/.test(v)) strength++;
            if (/[0-9]/.test(v) || /[^A-Za-z0-9]/.test(v)) strength++;

            if (strengthFill) {
                strengthFill.className = 'strength-fill';
                if (v.length === 0) {
                    strengthFill.className = 'strength-fill';
                } else if (strength <= 1) {
                    strengthFill.classList.add('weak');
                } else if (strength === 2) {
                    strengthFill.classList.add('medium');
                } else {
                    strengthFill.classList.add('strong');
                }
            }

            if (strengthText) {
                if (v.length === 0) {
                    strengthText.textContent = '';
                    strengthText.style.color = '';
                } else if (strength <= 1) {
                    strengthText.textContent = 'Faible';
                    strengthText.style.color = 'var(--danger)';
                } else if (strength === 2) {
                    strengthText.textContent = 'Moyen';
                    strengthText.style.color = 'var(--warning)';
                } else {
                    strengthText.textContent = 'Fort';
                    strengthText.style.color = 'var(--success)';
                }
            }

            if (v.length > 0 && v.length < 8) {
                showError('passwordError', 'Min. 8 caractères');
            } else {
                showError('passwordError', '');
            }

            // Re-validate confirm
            if (confirm && confirm.value) {
                if (confirm.value !== v) {
                    showError('confirmError', 'Les mots de passe ne correspondent pas');
                } else {
                    showError('confirmError', '');
                }
            }
        });
    }

    if (confirm) {
        confirm.addEventListener('input', function() {
            if (password && this.value !== password.value) {
                showError('confirmError', 'Les mots de passe ne correspondent pas');
            } else {
                showError('confirmError', '');
            }
        });
    }
})();

// ============================================
// AUTO-DISMISS ALERTS
// ============================================
(function() {
    document.querySelectorAll('.alert').forEach(function(el) {
        setTimeout(function() {
            el.style.transition = 'opacity 0.4s ease, transform 0.4s ease';
            el.style.opacity = '0';
            el.style.transform = 'translateY(-10px)';
            setTimeout(function() { el.remove(); }, 400);
        }, 4000);
    });
})();

// ============================================
// UTILITY
// ============================================
function escapeHtml(text) {
    var div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
