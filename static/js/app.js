/* LE VESTIAIRE — Forum Football JS */

// ============================================
// PITCH CANVAS BACKGROUND
// ============================================
(function drawPitch() {
    const canvas = document.getElementById('pitchCanvas');
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    let w, h;

    function resize() {
        w = canvas.width = window.innerWidth;
        h = canvas.height = window.innerHeight;
    }
    resize();
    window.addEventListener('resize', resize);

    function draw() {
        ctx.clearRect(0, 0, w, h);
        ctx.strokeStyle = '#22c55e';
        ctx.lineWidth = 1.5;

        const cx = w / 2, cy = h / 2;
        const pw = Math.min(w * 0.8, 900), ph = Math.min(h * 0.7, 600);
        const left = cx - pw / 2, top_ = cy - ph / 2;

        // Pitch outline
        ctx.strokeRect(left, top_, pw, ph);
        // Center line
        ctx.beginPath(); ctx.moveTo(cx, top_); ctx.lineTo(cx, top_ + ph); ctx.stroke();
        // Center circle
        ctx.beginPath(); ctx.arc(cx, cy, ph * 0.15, 0, Math.PI * 2); ctx.stroke();
        // Center dot
        ctx.beginPath(); ctx.arc(cx, cy, 3, 0, Math.PI * 2); ctx.fill();

        // Penalty areas
        const paW = pw * 0.16, paH = ph * 0.44;
        ctx.strokeRect(left, cy - paH / 2, paW, paH);
        ctx.strokeRect(left + pw - paW, cy - paH / 2, paW, paH);

        // Goal areas
        const gaW = pw * 0.06, gaH = ph * 0.2;
        ctx.strokeRect(left, cy - gaH / 2, gaW, gaH);
        ctx.strokeRect(left + pw - gaW, cy - gaH / 2, gaW, gaH);

        // Corner arcs
        const cr = pw * 0.02;
        [
            [left, top_, 0, Math.PI / 2],
            [left + pw, top_, Math.PI / 2, Math.PI],
            [left + pw, top_ + ph, Math.PI, Math.PI * 1.5],
            [left, top_ + ph, Math.PI * 1.5, Math.PI * 2]
        ].forEach(([x, y, s, e]) => {
            ctx.beginPath(); ctx.arc(x, y, cr, s, e); ctx.stroke();
        });
    }

    draw();
    window.addEventListener('resize', draw);
})();

// ============================================
// TOAST NOTIFICATION SYSTEM
// ============================================
const Toast = {
    container: null,
    init() { this.container = document.getElementById('toastContainer'); },
    show(message, type) {
        if (!this.container) this.init();
        if (!this.container) return;
        const el = document.createElement('div');
        el.className = 'toast toast-' + type;
        el.textContent = message;
        this.container.appendChild(el);
        setTimeout(() => {
            el.classList.add('toast-out');
            setTimeout(() => el.remove(), 300);
        }, 3000);
    },
    success(msg) { this.show(msg, 'success'); },
    error(msg) { this.show(msg, 'error'); },
    info(msg) { this.show(msg, 'info'); }
};

// ============================================
// CARD 3D TILT + LIGHT REFLECTION
// ============================================
function initCardTilt() {
    document.querySelectorAll('[data-tilt]').forEach(card => {
        card.addEventListener('mousemove', (e) => {
            const rect = card.getBoundingClientRect();
            const x = e.clientX - rect.left;
            const y = e.clientY - rect.top;
            const centerX = rect.width / 2;
            const centerY = rect.height / 2;
            const rotateX = ((y - centerY) / centerY) * -4;
            const rotateY = ((x - centerX) / centerX) * 4;

            card.style.transform = 'perspective(1000px) rotateX(' + rotateX + 'deg) rotateY(' + rotateY + 'deg) scale3d(1.01, 1.01, 1.01)';
            card.style.setProperty('--mouse-x', x + 'px');
            card.style.setProperty('--mouse-y', y + 'px');
        });

        card.addEventListener('mouseleave', () => {
            card.style.transform = '';
        });
    });
}

// ============================================
// LIKE PARTICLES
// ============================================
function spawnParticles(button, color) {
    const rect = button.getBoundingClientRect();
    const cx = rect.left + rect.width / 2;
    const cy = rect.top + rect.height / 2;
    const colors = color === 'red'
        ? ['#ef4444', '#f87171', '#fca5a5', '#22c55e', '#fbbf24']
        : ['#22c55e', '#4ade80', '#86efac', '#fbbf24', '#f59e0b'];

    for (let i = 0; i < 12; i++) {
        const p = document.createElement('div');
        p.className = 'particle';
        const size = Math.random() * 6 + 3;
        const angle = (Math.PI * 2 / 12) * i + Math.random() * 0.5;
        const dist = Math.random() * 50 + 30;
        p.style.width = size + 'px';
        p.style.height = size + 'px';
        p.style.left = cx + 'px';
        p.style.top = cy + 'px';
        p.style.background = colors[Math.floor(Math.random() * colors.length)];
        p.style.setProperty('--px', Math.cos(angle) * dist + 'px');
        p.style.setProperty('--py', Math.sin(angle) * dist + 'px');
        document.body.appendChild(p);
        setTimeout(() => p.remove(), 600);
    }
}

// ============================================
// SCOREBOARD COUNTER ANIMATION
// ============================================
function animateScoreboards() {
    document.querySelectorAll('.scoreboard').forEach(el => {
        const target = parseInt(el.dataset.target) || 0;
        if (target === 0) return;
        const duration = 800;
        const start = performance.now();
        const initial = 0;

        function tick(now) {
            const elapsed = now - start;
            const progress = Math.min(elapsed / duration, 1);
            const eased = 1 - Math.pow(1 - progress, 3);
            el.textContent = Math.round(initial + (target - initial) * eased);
            if (progress < 1) requestAnimationFrame(tick);
        }
        requestAnimationFrame(tick);
    });
}

// ============================================
// SORT BAR INDICATOR
// ============================================
function initSortIndicator() {
    const bar = document.getElementById('sortBar');
    if (!bar) return;
    const active = bar.querySelector('.sort-tab.active');
    const indicator = document.getElementById('sortIndicator');
    if (!active || !indicator) return;

    function position() {
        indicator.style.left = active.offsetLeft + 'px';
        indicator.style.width = active.offsetWidth + 'px';
    }
    position();
    window.addEventListener('resize', position);
}

// ============================================
// LIVE SEARCH
// ============================================
function initSearch() {
    const input = document.getElementById('searchInput');
    const dropdown = document.getElementById('searchDropdown');
    if (!input || !dropdown) return;

    let timer = null;

    input.addEventListener('input', () => {
        clearTimeout(timer);
        const q = input.value.trim();
        if (q.length < 2) { dropdown.classList.remove('active'); return; }

        timer = setTimeout(() => {
            fetch('/api/search?q=' + encodeURIComponent(q), {
                headers: { 'Accept': 'application/json' }
            })
            .then(r => r.json())
            .then(results => {
                if (!results || results.length === 0) {
                    dropdown.innerHTML = '<div class="search-no-result">Aucun résultat pour "' + q + '"</div>';
                } else {
                    dropdown.innerHTML = results.map(r =>
                        '<a href="/post/' + r.id + '" class="search-result">' +
                        '<div class="search-result-title">' + escapeHtml(r.title) + '</div>' +
                        '<div class="search-result-meta">' + escapeHtml(r.author) + ' &middot; ' + escapeHtml(r.excerpt) + '</div>' +
                        '</a>'
                    ).join('');
                }
                dropdown.classList.add('active');
            })
            .catch(() => { dropdown.classList.remove('active'); });
        }, 300);
    });

    input.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            const q = input.value.trim();
            if (q) window.location.href = '/search?q=' + encodeURIComponent(q);
        }
    });

    document.addEventListener('click', (e) => {
        if (!e.target.closest('#navSearchWrapper')) dropdown.classList.remove('active');
    });
}

// ============================================
// AJAX REACTIONS
// ============================================
function handleLike(btn, type) {
    const postId = btn.dataset.postId;
    fetch('/post/react/' + postId, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'Accept': 'application/json' },
        body: 'type=like'
    })
    .then(r => { if (r.status === 401) { window.location = '/login'; throw new Error('auth'); } return r.json(); })
    .then(data => {
        if (!data.success) return;
        // Update all like buttons for this post
        document.querySelectorAll('.like-btn[data-post-id="' + postId + '"]').forEach(b => {
            b.classList.toggle('active', data.userVote === 'like');
            const countEl = b.querySelector('.like-count');
            if (countEl) animateCount(countEl, data.likes);
        });
        document.querySelectorAll('.dislike-btn[data-post-id="' + postId + '"]').forEach(b => {
            b.classList.toggle('active', data.userVote === 'dislike');
            const countEl = b.querySelector('.dislike-count');
            if (countEl) animateCount(countEl, data.dislikes);
        });
        if (data.userVote === 'like') spawnParticles(btn, 'red');
    })
    .catch(e => { if (e.message !== 'auth') Toast.error('Erreur'); });
}

function handleDislike(btn, type) {
    const postId = btn.dataset.postId;
    fetch('/post/react/' + postId, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'Accept': 'application/json' },
        body: 'type=dislike'
    })
    .then(r => { if (r.status === 401) { window.location = '/login'; throw new Error('auth'); } return r.json(); })
    .then(data => {
        if (!data.success) return;
        document.querySelectorAll('.like-btn[data-post-id="' + postId + '"]').forEach(b => {
            b.classList.toggle('active', data.userVote === 'like');
            const countEl = b.querySelector('.like-count');
            if (countEl) animateCount(countEl, data.likes);
        });
        document.querySelectorAll('.dislike-btn[data-post-id="' + postId + '"]').forEach(b => {
            b.classList.toggle('active', data.userVote === 'dislike');
            const countEl = b.querySelector('.dislike-count');
            if (countEl) animateCount(countEl, data.dislikes);
        });
    })
    .catch(e => { if (e.message !== 'auth') Toast.error('Erreur'); });
}

function handleRepost(btn) {
    const postId = btn.dataset.postId;
    fetch('/repost/' + postId, {
        method: 'POST',
        headers: { 'Accept': 'application/json' }
    })
    .then(r => { if (r.status === 401) { window.location = '/login'; throw new Error('auth'); } return r.json(); })
    .then(data => {
        if (!data.success) return;
        document.querySelectorAll('.repost-btn[data-post-id="' + postId + '"]').forEach(b => {
            b.classList.toggle('active', data.reposted);
            const countEl = b.querySelector('.repost-count');
            if (countEl) animateCount(countEl, data.count);
        });
        if (data.reposted) { spawnParticles(btn, 'green'); Toast.success('Partagé !'); }
        else Toast.info('Partage retiré');
    })
    .catch(e => { if (e.message !== 'auth') Toast.error('Erreur'); });
}

function handleCommentReaction(btn, type) {
    const commentId = btn.dataset.commentId;
    const postId = btn.dataset.postId;
    fetch('/comment/react/' + commentId, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'Accept': 'application/json' },
        body: 'type=' + type + '&post_id=' + postId
    })
    .then(r => { if (r.status === 401) { window.location = '/login'; throw new Error('auth'); } return r.json(); })
    .then(data => {
        if (!data.success) return;
        const card = btn.closest('.comment-card') || btn.closest('.reply-card');
        if (!card) return;
        const likeBtn = card.querySelector('.like-btn[data-comment-id="' + commentId + '"]');
        const dislikeBtn = card.querySelector('.dislike-btn[data-comment-id="' + commentId + '"]');
        if (likeBtn) {
            likeBtn.classList.toggle('active', data.userVote === 'like');
            const c = likeBtn.querySelector('.like-count');
            if (c) animateCount(c, data.likes);
        }
        if (dislikeBtn) {
            dislikeBtn.classList.toggle('active', data.userVote === 'dislike');
            const c = dislikeBtn.querySelector('.dislike-count');
            if (c) animateCount(c, data.dislikes);
        }
        if (data.userVote === 'like') spawnParticles(btn, 'red');
    })
    .catch(e => { if (e.message !== 'auth') Toast.error('Erreur'); });
}

function animateCount(el, target) {
    const current = parseInt(el.textContent) || 0;
    if (current === target) return;
    el.style.transform = 'translateY(-100%)';
    el.style.opacity = '0';
    setTimeout(() => {
        el.textContent = target;
        el.style.transform = 'translateY(100%)';
        requestAnimationFrame(() => {
            el.style.transition = 'all 0.2s ease';
            el.style.transform = 'translateY(0)';
            el.style.opacity = '1';
            setTimeout(() => { el.style.transition = ''; }, 200);
        });
    }, 100);
}

// ============================================
// SHARE
// ============================================
function toggleShare(btn) {
    const menu = btn.closest('.share-wrapper').querySelector('.share-menu');
    document.querySelectorAll('.share-menu.active').forEach(m => { if (m !== menu) m.classList.remove('active'); });
    menu.classList.toggle('active');
}

function copyLink(path) {
    const url = window.location.origin + path;
    navigator.clipboard.writeText(url).then(() => Toast.success('Lien copié !')).catch(() => Toast.error('Erreur'));
    document.querySelectorAll('.share-menu.active').forEach(m => m.classList.remove('active'));
}

function shareTwitter(title, path) {
    const url = window.location.origin + path;
    window.open('https://twitter.com/intent/tweet?text=' + encodeURIComponent(title) + '&url=' + encodeURIComponent(url), '_blank');
    document.querySelectorAll('.share-menu.active').forEach(m => m.classList.remove('active'));
}

function shareWhatsApp(title, path) {
    const url = window.location.origin + path;
    window.open('https://wa.me/?text=' + encodeURIComponent(title + ' ' + url), '_blank');
    document.querySelectorAll('.share-menu.active').forEach(m => m.classList.remove('active'));
}

// Close share menus on outside click
document.addEventListener('click', (e) => {
    if (!e.target.closest('.share-wrapper')) {
        document.querySelectorAll('.share-menu.active').forEach(m => m.classList.remove('active'));
    }
});

// ============================================
// COMMENT FORM
// ============================================
function initCommentForm() {
    const textarea = document.getElementById('commentText');
    const submit = document.getElementById('commentSubmit');
    if (!textarea || !submit) return;

    textarea.addEventListener('input', () => {
        submit.disabled = textarea.value.trim().length === 0;
    });
}

// ============================================
// REPLY FORM
// ============================================
function toggleReplyForm(btn, commentId, postId) {
    const body = btn.closest('.comment-body');
    let existing = body.querySelector('.reply-form');
    if (existing) { existing.remove(); return; }

    const form = document.createElement('form');
    form.className = 'reply-form';
    form.action = '/comment/create';
    form.method = 'POST';
    form.innerHTML =
        '<input type="hidden" name="post_id" value="' + postId + '">' +
        '<input type="hidden" name="parent_id" value="' + commentId + '">' +
        '<textarea name="content" placeholder="Ta réponse..." rows="2" required></textarea>' +
        '<button type="submit" class="btn-primary btn-sm">Répondre</button>';
    body.appendChild(form);
    form.querySelector('textarea').focus();
}

// ============================================
// USER DROPDOWN
// ============================================
function initUserDropdown() {
    const navUser = document.getElementById('navUser');
    const dropdown = document.getElementById('userDropdown');
    if (!navUser || !dropdown) return;

    navUser.addEventListener('click', (e) => {
        e.stopPropagation();
        dropdown.classList.toggle('active');
    });

    document.addEventListener('click', (e) => {
        if (!e.target.closest('#navUser')) dropdown.classList.remove('active');
    });
}

// ============================================
// HAMBURGER MENU
// ============================================
function initHamburger() {
    const hamburger = document.getElementById('hamburger');
    const sidebar = document.getElementById('sidebarLeft');
    const overlay = document.getElementById('sidebarOverlay');
    if (!hamburger) return;

    hamburger.addEventListener('click', () => {
        hamburger.classList.toggle('active');
        if (sidebar) sidebar.classList.toggle('active');
        if (overlay) overlay.classList.toggle('active');
    });

    if (overlay) {
        overlay.addEventListener('click', () => {
            hamburger.classList.remove('active');
            sidebar.classList.remove('active');
            overlay.classList.remove('active');
        });
    }
}

// ============================================
// REGISTER STEPPER
// ============================================
function nextStep(step) {
    // Validate current step
    if (step === 2) {
        const username = document.getElementById('reg-username');
        const email = document.getElementById('reg-email');
        const usernameErr = document.getElementById('usernameError');
        const emailErr = document.getElementById('emailError');
        let valid = true;

        if (!username.value.trim() || username.value.trim().length < 3) {
            usernameErr.textContent = 'Minimum 3 caractères';
            valid = false;
        } else { usernameErr.textContent = ''; }

        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email.value.trim())) {
            emailErr.textContent = 'Email invalide';
            valid = false;
        } else { emailErr.textContent = ''; }

        if (!valid) return;
    }

    if (step === 3) {
        const password = document.getElementById('reg-password');
        const confirm = document.getElementById('reg-confirm');
        const passErr = document.getElementById('passwordError');
        const confErr = document.getElementById('confirmError');
        let valid = true;

        if (password.value.length < 6) {
            passErr.textContent = 'Minimum 6 caractères';
            valid = false;
        } else { passErr.textContent = ''; }

        if (password.value !== confirm.value) {
            confErr.textContent = 'Les mots de passe ne correspondent pas';
            valid = false;
        } else { confErr.textContent = ''; }

        if (!valid) return;
    }

    goToStep(step);
}

function prevStep(step) {
    goToStep(step);
}

function goToStep(step) {
    document.querySelectorAll('.form-step').forEach(s => s.classList.remove('active'));
    document.querySelectorAll('.step').forEach(s => {
        const sNum = parseInt(s.dataset.step);
        s.classList.remove('active', 'done');
        if (sNum < step) s.classList.add('done');
        if (sNum === step) s.classList.add('active');
    });
    const target = document.querySelector('.form-step[data-step="' + step + '"]');
    if (target) target.classList.add('active');
}

// ============================================
// PASSWORD STRENGTH
// ============================================
function initPasswordStrength() {
    const password = document.getElementById('reg-password');
    const fill = document.getElementById('strengthFill');
    const label = document.getElementById('strengthLabel');
    if (!password || !fill) return;

    password.addEventListener('input', () => {
        const val = password.value;
        let score = 0;
        if (val.length >= 6) score++;
        if (val.length >= 10) score++;
        if (/[A-Z]/.test(val)) score++;
        if (/[0-9]/.test(val)) score++;
        if (/[^A-Za-z0-9]/.test(val)) score++;

        const levels = [
            { w: '0%', c: '', t: '' },
            { w: '20%', c: '#ef4444', t: 'Très faible' },
            { w: '40%', c: '#f97316', t: 'Faible' },
            { w: '60%', c: '#eab308', t: 'Moyen' },
            { w: '80%', c: '#22c55e', t: 'Fort' },
            { w: '100%', c: '#16a34a', t: 'Très fort' }
        ];

        const level = levels[score] || levels[0];
        fill.style.width = level.w;
        fill.style.background = level.c;
        if (label) { label.textContent = level.t; label.style.color = level.c; }
    });
}

// ============================================
// IMAGE UPLOAD PREVIEW
// ============================================
function initFileUpload() {
    const upload = document.getElementById('fileUpload');
    if (!upload) return;
    const input = upload.querySelector('input[type="file"]');
    if (!input) return;

    input.addEventListener('change', () => {
        const file = input.files[0];
        if (!file) return;

        // Remove old preview if any
        const oldPreview = upload.querySelector('.upload-preview');
        if (oldPreview) oldPreview.remove();

        const reader = new FileReader();
        reader.onload = (e) => {
            const preview = document.createElement('div');
            preview.className = 'upload-preview';
            preview.innerHTML =
                '<img src="' + e.target.result + '" alt="Aperçu">' +
                '<button type="button" class="upload-remove" onclick="removeUpload(this)">&times;</button>' +
                '<span class="upload-filename">' + file.name + '</span>';
            upload.querySelector('.file-upload-label').style.display = 'none';
            upload.appendChild(preview);
        };
        reader.readAsDataURL(file);
    });
}

function removeUpload(btn) {
    const upload = btn.closest('.file-upload');
    const input = upload.querySelector('input[type="file"]');
    const preview = upload.querySelector('.upload-preview');
    const label = upload.querySelector('.file-upload-label');
    // Reset file input
    input.value = '';
    if (preview) preview.remove();
    if (label) label.style.display = '';
}

// ============================================
// TICKER DUPLICATION (for seamless loop)
// ============================================
function initTicker() {
    const track = document.getElementById('tickerTrack');
    if (!track) return;
    // Duplicate content for seamless loop
    track.innerHTML += track.innerHTML;
}

// ============================================
// UTILITY
// ============================================
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ============================================
// INIT
// ============================================
document.addEventListener('DOMContentLoaded', () => {
    initCardTilt();
    initSortIndicator();
    initSearch();
    initCommentForm();
    initUserDropdown();
    initHamburger();
    initPasswordStrength();
    initFileUpload();
    initTicker();
    animateScoreboards();

    // Close menus on Escape
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            document.querySelectorAll('.share-menu.active, .user-dropdown.active, .search-dropdown.active').forEach(m => m.classList.remove('active'));
        }
    });
});
