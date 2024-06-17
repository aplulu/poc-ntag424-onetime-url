package me.aplulu.ntag424writer

import android.content.Intent
import android.os.Bundle
import androidx.preference.Preference
import androidx.preference.PreferenceFragmentCompat

class SettingsFragment: PreferenceFragmentCompat() {
    override fun onCreatePreferences(savedInstanceState: Bundle?, rootKey: String?) {
        setPreferencesFromResource(R.xml.preferences, rootKey)

        val aboutPreference: Preference? = findPreference("about")
        aboutPreference?.onPreferenceClickListener = Preference.OnPreferenceClickListener {
            val intent = Intent(activity, AboutActivity::class.java)
            startActivity(intent)
            true
        }
    }
}