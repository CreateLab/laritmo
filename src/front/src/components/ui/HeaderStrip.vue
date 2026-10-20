<template>
  <div class="h-7 bg-sunken border-b border-border overflow-hidden" aria-hidden="true">
    <svg width="100%" height="28" class="block">
      <defs>
        <pattern id="strip" width="128" height="28" patternUnits="userSpaceOnUse">
          <g transform="translate(8 4) scale(0.8)">
            <g fill="var(--action)"><ellipse cx="12" cy="13" rx="9" ry="6"/><circle cx="7.5" cy="7.5" r="3.2"/><circle cx="16.5" cy="7.5" r="3.2"/><ellipse cx="4" cy="17.5" rx="3.5" ry="2"/><ellipse cx="20" cy="17.5" rx="3.5" ry="2"/></g>
            <circle cx="7.5" cy="7.3" r="1.9" fill="var(--sunken)"/><circle cx="16.5" cy="7.3" r="1.9" fill="var(--sunken)"/>
            <circle cx="7.9" cy="7.4" r="0.9" fill="var(--ink)"/><circle cx="16.9" cy="7.4" r="0.9" fill="var(--ink)"/>
          </g>
          <g transform="translate(46 4)">
            <path d="M2 11a8 7.5 0 0 1 16 0z" fill="var(--accent)"/><circle cx="7" cy="8" r="1.4" fill="var(--sunken)"/><circle cx="12.5" cy="6.8" r="1" fill="var(--sunken)"/>
            <rect x="7.5" y="11" width="5" height="6.5" rx="2" fill="var(--border-strong)"/>
          </g>
          <g transform="translate(76 8) scale(0.7)">
            <path d="M2 11a8 7.5 0 0 1 16 0z" fill="var(--accent)"/><circle cx="7" cy="8" r="1.4" fill="var(--sunken)"/><circle cx="12.5" cy="6.8" r="1" fill="var(--sunken)"/>
            <rect x="7.5" y="11" width="5" height="6.5" rx="2" fill="var(--border-strong)"/>
          </g>
          <g transform="translate(122 4) scale(-0.8 0.8)">
            <g fill="var(--action)"><ellipse cx="12" cy="13" rx="9" ry="6"/><circle cx="7.5" cy="7.5" r="3.2"/><circle cx="16.5" cy="7.5" r="3.2"/><ellipse cx="4" cy="17.5" rx="3.5" ry="2"/><ellipse cx="20" cy="17.5" rx="3.5" ry="2"/></g>
            <circle cx="7.5" cy="7.3" r="1.9" fill="var(--sunken)"/><circle cx="16.5" cy="7.3" r="1.9" fill="var(--sunken)"/>
            <circle cx="7.9" cy="7.4" r="0.9" fill="var(--ink)"/><circle cx="16.9" cy="7.4" r="0.9" fill="var(--ink)"/>
          </g>
        </pattern>
      </defs>
      <rect width="100%" height="28" fill="url(#strip)"/>
    </svg>
  </div>
</template>
