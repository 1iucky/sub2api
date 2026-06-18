export default {
  batchImageGuide: {
    title: 'Batch Image Generation',
    description: 'Submit multiple prompts in one job and download the generated images when complete'
  },
  // Home Page
  home: {
    viewOnGithub: 'View on GitHub',
    viewDocs: 'View Documentation',
    docs: 'Docs',
    switchToLight: 'Switch to Light Mode',
    switchToDark: 'Switch to Dark Mode',
    dashboard: 'Console',
    login: 'Sign in',
    getStarted: 'Request access',
    goToDashboard: 'Open console',
    eyebrow: 'AI API GATEWAY',
    // User-focused value proposition
    heroSubtitle: 'A single endpoint for every upstream AI model.',
    heroDescription: 'SiliconBase pools your Claude, OpenAI and Gemini subscriptions behind one API key — with authentication, usage billing and load balancing handled at the gateway.',
    // v2 homepage sections (factory.ai motion rebuild)
    nav: {
      gateway: 'Gateway',
      pools: 'Pools',
      docs: 'Docs',
      openMenu: 'Open menu',
      closeMenu: 'Close menu'
    },
    hero2: {
      eyebrow: 'AI API GATEWAY',
      title: 'One gateway for every AI subscription.',
      subhead:
        'Route Claude, OpenAI, Gemini, Bedrock and Antigravity through a single API key. Pool accounts, balance load, meter usage — at the gateway.',
      ctaPrimary: 'Get API key',
      ctaSecondary: 'Read the docs',
      demo: {
        title: 'gateway — live',
        sidebarGateway: 'Gateway',
        sidebarPools: 'Pools',
        sidebarKeys: 'Keys',
        sidebarSettings: 'Settings',
        kpiRpm: 'RPM',
        kpiRpmValue: '1.2k',
        kpiP95: 'P95',
        kpiP95Value: '240ms',
        kpiUptime: 'Uptime',
        kpiUptimeValue: '99.9%',
        kpiKeys: 'Keys',
        kpiKeysValue: '348',
        health: 'health',
        streamTitle: 'request stream',
        streamRoute: 'route',
        streamUpstream: 'upstream',
        streamStatus: '200 ok',
        sparkClaude: 'claude',
        sparkOpenai: 'openai',
        sparkGemini: 'gemini'
      }
    },
    marquee2: {
      label: 'Compatible with',
      items: 'Claude · OpenAI · Gemini · Bedrock · Antigravity'
    },
    bento2: {
      eyebrow: 'WHAT THE GATEWAY DOES',
      title: 'Defining your gateway',
      cards: {
        pooling: {
          index: '01',
          title: 'Multi-upstream pooling',
          desc: 'Aggregate Claude, OpenAI and Gemini accounts behind one key. Sticky sessions, per-account RPM and concurrency limits, automatic failover.'
        },
        routing: {
          index: '02',
          title: 'Sticky sessions & load balancing',
          desc: 'Pin a session to an account, distribute traffic across the pool, and keep latency flat as you scale keys and teams.'
        },
        billing: {
          index: '03',
          title: 'Quotas & metered billing',
          desc: 'Per-key usage is metered in real time with configurable quotas and spend limits. Team consumption stays observable and bounded.'
        }
      }
    },
    cta2: {
      eyebrow: 'SHIP WITHOUT MANAGING QUOTAS',
      title: 'Ship AI features without managing quotas.',
      button: 'Start free',
      note: 'No credit card to try. Pool your first upstream in minutes.'
    },
    footer2: {
      tagline: 'Operated as SiliconBase',
      domain: 'siliconbase.link',
      columns: {
        product: 'Product',
        resources: 'Resources',
        legal: 'Legal'
      },
      links: {
        gateway: 'Gateway',
        pools: 'Pools',
        pricing: 'Pricing',
        docs: 'Docs',
        status: 'Status',
        changelog: 'Changelog',
        privacy: 'Privacy',
        terms: 'Terms'
      },
      copyright: '© 2026 SiliconBase. All rights reserved.'
    },
    tags: {
      subscriptionToApi: 'Subscription to API',
      stickySession: 'Session affinity',
      realtimeBilling: 'Metered billing'
    },
    // Pain points section
    painPoints: {
      title: 'Sound Familiar?',
      items: {
        expensive: {
          title: 'High Subscription Costs',
          desc: 'Paying for multiple AI subscriptions that add up every month'
        },
        complex: {
          title: 'Account Chaos',
          desc: 'Managing scattered accounts and API keys across different platforms'
        },
        unstable: {
          title: 'Service Interruptions',
          desc: 'Single accounts hitting rate limits and disrupting your workflow'
        },
        noControl: {
          title: 'No Usage Control',
          desc: "Can't track where your money goes or limit team member usage"
        }
      }
    },
    // Solutions section
    solutions: {
      title: 'We Solve These Problems',
      subtitle: 'Three simple steps to stress-free AI access'
    },
    features: {
      unifiedGateway: 'Unified gateway',
      unifiedGatewayDesc: 'One API key, one endpoint. Requests are routed to the correct upstream model — Claude, GPT, Gemini — based on key group, transparent to the client.',
      multiAccount: 'Account pooling & failover',
      multiAccountDesc: 'Multiple upstream accounts are balanced per group with sticky sessions, RPM and concurrency limits, and automatic failover when a session errors.',
      balanceQuota: 'Metered billing & quotas',
      balanceQuotaDesc: 'Per-key usage is metered in real time with configurable quotas and spend limits, so team consumption stays observable and bounded.'
    },
    // Comparison section
    comparison: {
      title: 'Why Choose Us?',
      headers: {
        feature: 'Comparison',
        official: 'Official Subscriptions',
        us: 'Our Platform'
      },
      items: {
        pricing: {
          feature: 'Pricing',
          official: 'Fixed monthly fee, pay even if unused',
          us: 'Pay only for what you use'
        },
        models: {
          feature: 'Model Selection',
          official: 'Single provider only',
          us: 'Switch between models freely'
        },
        management: {
          feature: 'Account Management',
          official: 'Manage each service separately',
          us: 'Unified key, one dashboard'
        },
        stability: {
          feature: 'Stability',
          official: 'Single account rate limits',
          us: 'Multi-account pool, auto-failover'
        },
        control: {
          feature: 'Usage Control',
          official: 'Not available',
          us: 'Quotas & detailed analytics'
        }
      }
    },
    providers: {
      title: 'Upstream providers',
      description: 'One endpoint, several upstreams — routed per key group.',
      supported: 'Routed',
      soon: 'Planned',
      claude: 'Claude',
      gemini: 'Gemini',
      antigravity: 'Antigravity',
      more: 'More'
    },
    // CTA section
    cta: {
      title: 'Ready to Get Started?',
      description: 'Sign up now and get free trial credits to experience seamless AI access',
      button: 'Sign Up Free'
    },
    footer: {
      allRightsReserved: 'All rights reserved.',
      tagline: 'Operated as SiliconBase'
    }
  },

  // Key Usage Query Page
  keyUsage: {
    title: 'API Key Usage',
    subtitle: 'Enter your API Key to view real-time spending and usage status',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: 'Query',
    querying: 'Querying...',
    privacyNote: 'Your Key is processed locally in the browser and will not be stored',
    dateRange: 'Date Range:',
    dateRangeToday: 'Today',
    dateRange7d: '7 Days',
    dateRange30d: '30 Days',
    dateRange90d: '90 Days',
    dateRangeCustom: 'Custom',
    apply: 'Apply',
    used: 'Used',
    detailInfo: 'Detail Information',
    tokenStats: 'Token Statistics',
    dailyDetail: 'Daily Detail',
    modelStats: 'Model Usage Statistics',
    // Table headers
    date: 'Date',
    model: 'Model',
    requests: 'Requests',
    inputTokens: 'Input Tokens',
    outputTokens: 'Output Tokens',
    cacheCreationTokens: 'Cache Creation',
    cacheReadTokens: 'Cache Read',
    cacheWriteTokens: 'Cache Write',
    totalTokens: 'Total Tokens',
    cost: 'Cost',
    // Status
    quotaMode: 'Key Quota Mode',
    walletBalance: 'Wallet Balance',
    // Ring card titles
    totalQuota: 'Total Quota',
    limit5h: '5-Hour Limit',
    limitDaily: 'Daily Limit',
    limit7d: '7-Day Limit',
    limitWeekly: 'Weekly Limit',
    limitMonthly: 'Monthly Limit',
    // Detail rows
    remainingQuota: 'Remaining Quota',
    expiresAt: 'Expires At',
    todayExpires: '(expires today)',
    daysLeft: '({days} days)',
    usedQuota: 'Used Quota',
    resetNow: 'Resetting soon',
    subscriptionType: 'Subscription Type',
    subscriptionExpires: 'Subscription Expires',
    // Usage stat cells
    todayRequests: 'Today Requests',
    todayInputTokens: 'Today Input',
    todayOutputTokens: 'Today Output',
    todayTokens: 'Today Tokens',
    todayCacheCreation: 'Today Cache Creation',
    todayCacheRead: 'Today Cache Read',
    todayCost: 'Today Cost',
    rpmTpm: 'RPM / TPM',
    totalRequests: 'Total Requests',
    totalInputTokens: 'Total Input',
    totalOutputTokens: 'Total Output',
    totalTokensLabel: 'Total Tokens',
    totalCacheCreation: 'Total Cache Creation',
    totalCacheRead: 'Total Cache Read',
    totalCost: 'Total Cost',
    avgDuration: 'Avg Duration',
    // Messages
    enterApiKey: 'Please enter an API Key',
    querySuccess: 'Query successful',
    queryFailed: 'Query failed',
    queryFailedRetry: 'Query failed, please try again later',
    noDailyUsage: 'No daily usage data',
  },

  // Setup Wizard
  setup: {
    title: 'SiliconBase Setup',
    description: 'Configure your SiliconBase instance',
    database: {
      title: 'Database Configuration',
      description: 'Connect to your PostgreSQL database',
      host: 'Host',
      port: 'Port',
      username: 'Username',
      password: 'Password',
      databaseName: 'Database Name',
      sslMode: 'SSL Mode',
      passwordPlaceholder: 'Password',
      ssl: {
        disable: 'Disable',
        require: 'Require',
        verifyCa: 'Verify CA',
        verifyFull: 'Verify Full'
      }
    },
    redis: {
      title: 'Redis Configuration',
      description: 'Connect to your Redis server',
      host: 'Host',
      port: 'Port',
      password: 'Password (optional)',
      database: 'Database',
      passwordPlaceholder: 'Password',
      enableTls: 'Enable TLS',
      enableTlsHint: 'Use TLS when connecting to Redis (public CA certs)'
    },
    admin: {
      title: 'Admin Account',
      description: 'Create your administrator account',
      email: 'Email',
      password: 'Password',
      confirmPassword: 'Confirm Password',
      passwordPlaceholder: 'Min 8 characters',
      confirmPasswordPlaceholder: 'Confirm password',
      passwordMismatch: 'Passwords do not match'
    },
    ready: {
      title: 'Ready to Install',
      description: 'Review your configuration and complete setup',
      database: 'Database',
      redis: 'Redis',
      adminEmail: 'Admin Email'
    },
    status: {
      testing: 'Testing...',
      success: 'Connection Successful',
      testConnection: 'Test Connection',
      installing: 'Installing...',
      completeInstallation: 'Complete Installation',
      completed: 'Installation completed!',
      redirecting: 'Redirecting to login page...',
      restarting: 'Service is restarting, please wait...',
      timeout: 'Service restart is taking longer than expected. Please refresh the page manually.'
    }
  },

  // Common
}
