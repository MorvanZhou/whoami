<script setup lang="ts">
const { t, locale } = useI18n()

const title = computed(() => t('seo.terms.title'))
const description = computed(() => t('seo.terms.description'))

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
      <h1>Terms of Service</h1>
      <p class="text-th-text-t text-sm">Effective date: {{ effectiveDate }}</p>

      <p>
        Welcome to <strong>whoami</strong> ("we", "us", or "our"), operated at
        <a href="https://whoamiagent.com">whoamiagent.com</a>. By accessing or
        using our website and services, you agree to be bound by these Terms of
        Service ("Terms"). If you do not agree, please do not use the service.
      </p>

      <h2>1. Service Description</h2>
      <p>
        whoami provides a persistent identity profile that syncs across
        AI agents including Cursor, Claude, Codex and other MCP-compatible
        tools. Users can create, edit and manage a personal profile that is
        read by connected AI agents to provide personalised context.
      </p>

      <h2>2. Accounts &amp; Authentication</h2>
      <p>
        You may sign in through third-party OAuth providers (e.g. GitHub,
        Google, Microsoft, Apple, X). We do not store your passwords. You
        are responsible for keeping your account credentials and connected
        agent keys secure.
      </p>

      <h2>3. User Content</h2>
      <p>
        You retain ownership of all content you submit to your identity
        profile ("User Content"). By using the service, you grant us a
        limited, non-exclusive licence to store, process and transmit your
        User Content solely for the purpose of providing and improving the
        service.
      </p>
      <p>
        You are solely responsible for the accuracy and legality of your
        User Content. You must not submit content that is unlawful,
        defamatory, infringing or harmful.
      </p>

      <h2>4. Acceptable Use</h2>
      <p>You agree not to:</p>
      <ul>
        <li>Use the service for any illegal or unauthorised purpose.</li>
        <li>Attempt to reverse-engineer, decompile or disassemble any part of the service.</li>
        <li>Interfere with or disrupt the service or its infrastructure.</li>
        <li>Use automated means to scrape, crawl or collect data from the service beyond normal API usage.</li>
        <li>Impersonate any person or entity.</li>
      </ul>

      <h2>5. Intellectual Property</h2>
      <p>
        All trademarks, logos and service marks displayed on the site are
        our property or the property of their respective owners. The source
        code is available under its published open-source licence; these
        Terms govern your use of the hosted service at whoamiagent.com.
      </p>

      <h2>6. Third-Party Services</h2>
      <p>
        Our service integrates with third-party AI tools and OAuth
        providers. We are not responsible for the practices or policies of
        these third parties. Your use of third-party services is subject to
        their respective terms and privacy policies.
      </p>

      <h2>7. Availability &amp; Modifications</h2>
      <p>
        We strive to keep the service available but do not guarantee
        uninterrupted access. We reserve the right to modify, suspend or
        discontinue the service (or any part of it) at any time, with or
        without notice.
      </p>

      <h2>8. Disclaimer of Warranties</h2>
      <p>
        The service is provided <strong>"as is"</strong> and
        <strong>"as available"</strong> without warranties of any kind,
        express or implied, including but not limited to merchantability,
        fitness for a particular purpose and non-infringement.
      </p>

      <h2>9. Limitation of Liability</h2>
      <p>
        To the maximum extent permitted by law, we shall not be liable for
        any indirect, incidental, special, consequential or punitive
        damages, or any loss of data, use, goodwill or profits, arising
        from your use of the service.
      </p>

      <h2>10. Termination</h2>
      <p>
        We may suspend or terminate your access at any time for conduct
        that we believe violates these Terms or is harmful to other users
        or the service. You may stop using the service at any time.
      </p>

      <h2>11. Changes to These Terms</h2>
      <p>
        We may update these Terms from time to time. When we do, we will
        revise the "Effective date" above. Continued use of the service
        after changes constitutes acceptance of the updated Terms.
      </p>

      <h2>12. Contact</h2>
      <p>
        If you have questions about these Terms, please open an issue on
        our <a href="https://github.com/MorvanZhou/whoami" target="_blank" rel="noopener noreferrer">GitHub repository</a>.
      </p>
    </article>

    <!-- 中文 -->
    <article v-else class="prose prose-neutral dark:prose-invert max-w-none">
      <h1>服务条款</h1>
      <p class="text-th-text-t text-sm">生效日期：{{ effectiveDate }}</p>

      <p>
        欢迎使用 <strong>whoami</strong>（以下简称"我们"），运营于
        <a href="https://whoamiagent.com">whoamiagent.com</a>。访问或使用我们的网站及服务即表示您同意受本服务条款（以下简称"条款"）的约束。如果您不同意，请停止使用本服务。
      </p>

      <h2>1. 服务说明</h2>
      <p>
        whoami 提供一个可跨 AI 工具同步的持久身份档案服务，支持 Cursor、Claude、Codex
        及其他 MCP 兼容工具。用户可以创建、编辑和管理个人档案，供已连接的 AI
        agent 读取，以提供个性化上下文。
      </p>

      <h2>2. 账号与认证</h2>
      <p>
        您可以通过第三方 OAuth 提供商（如 GitHub、Google、Microsoft、Apple、X）登录。我们不存储您的密码。您有责任妥善保管您的账号凭证和已连接的
        Agent 密钥。
      </p>

      <h2>3. 用户内容</h2>
      <p>
        您保留您提交到身份档案中的所有内容（以下简称"用户内容"）的所有权。使用本服务即表示您授予我们有限的、非排他性的许可，仅用于存储、处理和传输您的用户内容，以提供和改进服务。
      </p>
      <p>
        您对用户内容的准确性和合法性负全部责任。您不得提交任何违法、诽谤、侵权或有害的内容。
      </p>

      <h2>4. 可接受的使用方式</h2>
      <p>您同意不会：</p>
      <ul>
        <li>将本服务用于任何非法或未经授权的目的。</li>
        <li>尝试对本服务的任何部分进行逆向工程、反编译或反汇编。</li>
        <li>干扰或破坏本服务及其基础设施。</li>
        <li>使用超出正常 API 用量的自动化手段抓取、爬取或收集本服务上的数据。</li>
        <li>冒充任何个人或实体。</li>
      </ul>

      <h2>5. 知识产权</h2>
      <p>
        本网站上展示的所有商标、标识和服务标志均为我们或各自所有者的财产。源代码按其发布的开源许可证提供；本条款约束您对
        whoamiagent.com 上托管服务的使用。
      </p>

      <h2>6. 第三方服务</h2>
      <p>
        我们的服务集成了第三方 AI 工具和 OAuth
        提供商。我们不对这些第三方的做法或政策负责。您使用第三方服务须遵守其各自的条款和隐私政策。
      </p>

      <h2>7. 可用性与变更</h2>
      <p>
        我们努力保持服务的可用性，但不保证不间断的访问。我们保留随时修改、暂停或终止服务（或其任何部分）的权利，无论是否事先通知。
      </p>

      <h2>8. 免责声明</h2>
      <p>
        本服务按<strong>"原样"</strong>和<strong>"可用"</strong>的基础提供，不附带任何形式的明示或暗示保证，包括但不限于适销性、特定用途适用性和非侵权性。
      </p>

      <h2>9. 责任限制</h2>
      <p>
        在法律允许的最大范围内，我们不对因您使用本服务而产生的任何间接、附带、特殊、后果性或惩罚性损害，或任何数据丢失、使用损失、商誉损失或利润损失承担责任。
      </p>

      <h2>10. 终止</h2>
      <p>
        如果我们认为您的行为违反本条款或对其他用户或服务有害，我们可以随时暂停或终止您的访问权限。您可以随时停止使用本服务。
      </p>

      <h2>11. 条款变更</h2>
      <p>
        我们可能会不时更新本条款。届时，我们将修改上述"生效日期"。在变更后继续使用本服务即表示您接受更新后的条款。
      </p>

      <h2>12. 联系方式</h2>
      <p>
        如果您对本条款有任何疑问，请在我们的
        <a href="https://github.com/MorvanZhou/whoami" target="_blank" rel="noopener noreferrer">GitHub 仓库</a>
        提交 Issue。
      </p>
    </article>
  </div>
</template>
