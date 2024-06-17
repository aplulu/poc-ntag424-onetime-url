package me.aplulu.ntag424writer

import android.os.Bundle
import android.webkit.WebView
import androidx.appcompat.app.AppCompatActivity

class AboutActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_about)

        supportActionBar?.setTitle(R.string.about)

        supportActionBar?.setDisplayHomeAsUpEnabled(true)

        val webView: WebView = findViewById(R.id.webview)
        webView.loadUrl("file:///android_asset/about.html")
    }

    override fun onSupportNavigateUp(): Boolean {
        finish()
        return true
    }
}