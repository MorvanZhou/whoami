<script setup lang="ts">
const { t, locale } = useI18n()

const title = computed(() => t('seo.privacy.title'))
const description = computed(() => t('seo.privacy.description'))

useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description,
})

const effectiveDate = '2025-06-01'
</script>

<template>
  <div class="mx-auto max-w-3xl px-6 py-16 sm:py-24">
    <!-- English -->
    <article v-if="locale === 'en'" class="prose prose-neutral dark:prose-invert max-w-none">
      <h1>Privacy Policy</h1>
      <p class="text-th-text-t text-sm">Effective date: {{ effectiveDate }}</p>

      <p>
        This Privacy Policy explains how <strong>whoami</strong> ("we", "us",
        or "our"), operated at
        <a href="https://whoamiagent.com">whoamiagent.com</a>, collects, uses
        and protects your information. By using our service, you agree to the
        practices described below.
      </p>

      <h2>1. Information We Collect</h2>

      <h3>1.1 Account Information</h3>
      <p>
        When you sign in through a third-party OAuth provider (GitHub, Google,
        Microsoft, Apple, or X), we receive and store:
      </p>
      <ul>
        <li>Your display name</li>
        <li>Email address (if provided by the OAuth provider)</li>
        <li>Avatar / profile picture URL</li>
        <li>A unique provider-specific user ID</li>
      </ul>
      <p>We do <strong>not</strong> receive or store your passwords.</p>

      <h3>1.2 Identity Profile Content</h3>
      <p>
        You may create and edit a personal identity profile in Markdown
        format. This content is stored on our servers and made available to
        AI agents you have authorised via agent keys.
      </p>

      <h3>1.3 Usage Data</h3>
      <p>
        We use Google Analytics 4 to collect anonymised usage data such as
        page views, session duration, device type and referring sources. This
        data does not personally identify you.
      </p>

      <h2>2. How We Use Your Information</h2>
      <ul>
        <li><strong>Provide the service</strong> — authenticate you, store your profile and sync it to your authorised AI agents.</li>
        <li><strong>Improve the service</strong> — analyse aggregated usage patterns to fix bugs and prioritise features.</li>
        <li><strong>Communicate</strong> — send essential service-related notices (e.g. security alerts). We do not send marketing emails.</li>
      </ul>

      <h2>3. Data Sharing</h2>
      <p>We do <strong>not</strong> sell your personal information. We share data only in these cases:</p>
      <ul>
        <li><strong>With your AI agents</strong> — your identity profile is shared with AI agents you explicitly connect via agent keys.</li>
        <li><strong>Service providers</strong> — we use third-party infrastructure (cloud hosting, analytics) that may process data on our behalf under strict confidentiality obligations.</li>
        <li><strong>Legal requirements</strong> — we may disclose information if required by law, regulation or legal process.</li>
      </ul>

      <h2>4. Data Storage &amp; Security</h2>
      <p>
        Your data is stored on servers with industry-standard security
        measures including encrypted connections (HTTPS/TLS), secure
        authentication tokens (httpOnly cookies) and access controls. While
        we take reasonable precautions, no system is 100% secure.
      </p>

      <h2>5. Data Retention</h2>
      <p>
        We retain your account and profile data for as long as your account
        is active. If you delete your account, we will remove your personal
        data within 30 days, except where retention is required by law.
      </p>

      <h2>6. Cookies</h2>
      <p>We use the following cookies:</p>
      <ul>
        <li><strong>access_token</strong> — an httpOnly, secure session cookie for authentication. Essential for the service to function.</li>
        <li><strong>Google Analytics cookies</strong> — anonymised analytics cookies (<code>_ga</code>, <code>_ga_*</code>) to understand usage patterns.</li>
      </ul>

      <h2>7. Your Rights</h2>
      <p>You have the right to:</p>
      <ul>
        <li><strong>Access</strong> your personal data via the dashboard.</li>
        <li><strong>Edit or delete</strong> your identity profile at any time.</li>
        <li><strong>Delete your account</strong> by contacting us.</li>
        <li><strong>Withdraw consent</strong> for analytics by using browser-level cookie controls or ad-blockers.</li>
      </ul>

      <h2>8. Children's Privacy</h2>
      <p>
        Our service is not directed to children under 13. We do not
        knowingly collect personal information from children. If you believe
        a child has provided us with personal data, please contact us so we
        can remove it.
      </p>

      <h2>9. International Data Transfers</h2>
      <p>
        Your data may be processed in countries other than your own. By
        using the service, you consent to the transfer of your information
        to countries that may have different data protection laws.
      </p>

      <h2>10. Changes to This Policy</h2>
      <p>
        We may update this Privacy Policy from time to time. When we do, we
        will revise the "Effective date" above. Continued use of the service
        after changes constitutes acceptance of the updated policy.
      </p>

      <h2>11. Contact</h2>
      <p>
        If you have questions about this Privacy Policy, please open an
        issue on our
        <a href="https://github.com/MorvanZhou/whoami" target="_blank" rel="noopener noreferrer">GitHub repository</a>.
      </p>
    </article>

    <!-- 中文 -->
    <article v-else class="prose prose-neutral dark:prose-invert max-w-none">
      <h1>隐私政策</h1>
      <p class="text-th-text-t text-sm">生效日期：{{ effectiveDate }}</p>

      <p>
        本隐私政策说明 <strong>whoami</strong>（以下简称"我们"）如何收集、使用和保护您的信息。whoami 运营于
        <a href="https://whoamiagent.com">whoamiagent.com</a>。使用我们的服务即表示您同意以下做法。
      </p>

      <h2>1. 我们收集的信息</h2>

      <h3>1.1 账号信息</h3>
      <p>
        当您通过第三方 OAuth 提供商（GitHub、Google、Microsoft、Apple 或 X）登录时，我们会接收并存储：
      </p>
      <ul>
        <li>您的显示名称</li>
        <li>电子邮箱地址（如果 OAuth 提供商提供）</li>
        <li>头像 / 个人资料图片 URL</li>
        <li>提供商特定的唯一用户 ID</li>
      </ul>
      <p>我们<strong>不会</strong>接收或存储您的密码。</p>

      <h3>1.2 身份档案内容</h3>
      <p>
        您可以用 Markdown 格式创建和编辑个人身份档案。该内容存储在我们的服务器上，并提供给您通过
        Agent 密钥授权的 AI agent 访问。
      </p>

      <h3>1.3 使用数据</h3>
      <p>
        我们使用 Google Analytics 4 收集匿名使用数据，例如页面浏览量、会话时长、设备类型和来源。这些数据不会识别您的个人身份。
      </p>

      <h2>2. 我们如何使用您的信息</h2>
      <ul>
        <li><strong>提供服务</strong> — 验证您的身份、存储档案并同步到您授权的 AI agent。</li>
        <li><strong>改进服务</strong> — 分析聚合使用模式，以修复问题和确定功能优先级。</li>
        <li><strong>通信</strong> — 发送必要的服务相关通知（如安全警报）。我们不发送营销邮件。</li>
      </ul>

      <h2>3. 数据共享</h2>
      <p>我们<strong>不会</strong>出售您的个人信息。仅在以下情况下共享数据：</p>
      <ul>
        <li><strong>与您的 AI agent</strong> — 您的身份档案会共享给您通过 Agent 密钥明确连接的 AI agent。</li>
        <li><strong>服务提供商</strong> — 我们使用第三方基础设施（云托管、分析），它们可能在严格的保密义务下代表我们处理数据。</li>
        <li><strong>法律要求</strong> — 如果法律、法规或法律程序要求，我们可能会披露信息。</li>
      </ul>

      <h2>4. 数据存储与安全</h2>
      <p>
        您的数据存储在采用行业标准安全措施的服务器上，包括加密连接（HTTPS/TLS）、安全认证令牌（httpOnly
        cookie）和访问控制。虽然我们采取了合理的预防措施，但没有任何系统能保证 100% 安全。
      </p>

      <h2>5. 数据保留</h2>
      <p>
        只要您的账号处于活跃状态，我们就会保留您的账号和档案数据。如果您删除账号，我们将在 30
        天内删除您的个人数据，法律要求保留的除外。
      </p>

      <h2>6. Cookie</h2>
      <p>我们使用以下 Cookie：</p>
      <ul>
        <li><strong>access_token</strong> — 用于身份验证的 httpOnly 安全会话 Cookie，是服务正常运行所必需的。</li>
        <li><strong>Google Analytics Cookie</strong> — 匿名分析 Cookie（<code>_ga</code>、<code>_ga_*</code>），用于了解使用模式。</li>
      </ul>

      <h2>7. 您的权利</h2>
      <p>您有权：</p>
      <ul>
        <li>通过控制台<strong>访问</strong>您的个人数据。</li>
        <li>随时<strong>编辑或删除</strong>您的身份档案。</li>
        <li>通过联系我们<strong>删除您的账号</strong>。</li>
        <li>使用浏览器级别的 Cookie 控制或广告拦截器<strong>撤回分析同意</strong>。</li>
      </ul>

      <h2>8. 儿童隐私</h2>
      <p>
        我们的服务不面向 13 岁以下的儿童。我们不会故意收集儿童的个人信息。如果您认为有儿童向我们提供了个人数据，请联系我们以便我们删除。
      </p>

      <h2>9. 国际数据传输</h2>
      <p>
        您的数据可能在您所在国家/地区以外的地方处理。使用本服务即表示您同意将您的信息传输到可能具有不同数据保护法律的国家/地区。
      </p>

      <h2>10. 政策变更</h2>
      <p>
        我们可能会不时更新本隐私政策。届时，我们将修改上述"生效日期"。在变更后继续使用本服务即表示您接受更新后的政策。
      </p>

      <h2>11. 联系方式</h2>
      <p>
        如果您对本隐私政策有任何疑问，请在我们的
        <a href="https://github.com/MorvanZhou/whoami" target="_blank" rel="noopener noreferrer">GitHub 仓库</a>
        提交 Issue。
      </p>
    </article>
  </div>
</template>
