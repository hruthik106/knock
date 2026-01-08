<h1>knock</h1>

<p>
  <strong>knock</strong> is a lightweight command-line tool for checking the health
  of HTTP endpoints.
</p>

<p>
  It reports whether an endpoint is <em>alive</em>, <em>unhealthy</em>, or
  <em>unreachable</em>, and is designed to be fast, minimal, and script-friendly.
</p>

<hr />

<h2>Installation</h2>

<p><strong>Using Go:</strong></p>

<pre>
go install github.com/hruthik106/knock/cmd/knock@v0.1.1
</pre>

<p>
  Ensure <code>$HOME/go/bin</code> (Linux/macOS) or
  <code>%USERPROFILE%\go\bin</code> (Windows) is in your <code>PATH</code>.
</p>

<hr />

<h2>Usage</h2>

<h3>Check a single URL</h3>

<pre>
knock &lt;url&gt;
</pre>

<h3>Check multiple URLs from a file</h3>

<p>
  The file must contain one URL per line.
</p>

<pre>
knock -f &lt;file&gt;
</pre>

<hr />

<h2>Flags</h2>

<table>
  <thead>
    <tr>
      <th>Flag</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>-f</code>, <code>--file &lt;file&gt;</code></td>
      <td>Read targets from a file (one URL per line).</td>
    </tr>
    <tr>
      <td><code>-t</code>, <code>--timeout &lt;duration&gt;</code></td>
      <td>
        Request timeout (Go duration format).<br />
        Default: <code>5s</code>
      </td>
    </tr>
    <tr>
      <td><code>--method &lt;HEAD|GET&gt;</code></td>
      <td>
        HTTP method to use for the request.<br />
        Default: <code>HEAD</code>
      </td>
    </tr>
    <tr>
      <td><code>-o</code>, <code>--only &lt;filter&gt;</code></td>
      <td>
        Filter output by result type.<br />
        Allowed values:
        <code>al</code>, <code>alive</code>,
        <code>uh</code>, <code>unhealthy</code>,
        <code>ur</code>, <code>unreachable</code>
      </td>
    </tr>
  </tbody>
</table>

<hr />

<h2>Output Filtering</h2>

<p>
  The <code>--only</code> flag controls which results are printed.
  All targets are still checked; filtering affects output only.
</p>

<pre>
knock -o al -f &lt;file&gt;   <!-- show only alive targets -->
knock -o uh -f &lt;file&gt;   <!-- show only unhealthy targets -->
knock -o ur -f &lt;file&gt;   <!-- show only unreachable targets -->
</pre>

<hr />

<h2>Exit Codes</h2>

<table>
  <thead>
    <tr>
      <th>Code</th>
      <th>Meaning</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>0</code></td>
      <td>All targets are alive</td>
    </tr>
    <tr>
      <td><code>1</code></td>
      <td>At least one target is unhealthy</td>
    </tr>
    <tr>
      <td><code>3</code></td>
      <td>At least one target is unreachable</td>
    </tr>
    <tr>
      <td><code>2</code></td>
      <td>Usage or flag error</td>
    </tr>
  </tbody>
</table>

<hr />

<h2>Status</h2>

<p>
  <strong>knock</strong> is in early development (<code>v0.1.1</code>).
  The core behavior is stable, but new features may be added.
</p>
