package me.aplulu.ntag424writer

import android.app.PendingIntent
import android.content.Intent
import android.media.MediaPlayer
import android.nfc.NfcAdapter
import android.nfc.Tag
import android.nfc.tech.IsoDep
import android.os.Bundle
import android.util.Log
import android.view.Menu
import android.view.MenuItem
import android.widget.Toast
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.appcompat.widget.Toolbar
import androidx.core.content.IntentCompat
import net.bplearning.ntag424.DnaCommunicator
import net.bplearning.ntag424.card.KeyInfo
import net.bplearning.ntag424.card.KeySet
import net.bplearning.ntag424.command.ChangeFileSettings
import net.bplearning.ntag424.command.GetCardUid
import net.bplearning.ntag424.command.GetFileSettings
import net.bplearning.ntag424.command.WriteData
import net.bplearning.ntag424.constants.Ntag424
import net.bplearning.ntag424.constants.Permissions
import net.bplearning.ntag424.encryptionmode.AESEncryptionMode
import net.bplearning.ntag424.sdm.NdefTemplateMaster
import net.bplearning.ntag424.sdm.SDMSettings
import net.bplearning.ntag424.util.ByteUtil

class MainActivity : AppCompatActivity() {
    companion object {
        private const val TAG = "NTAG424Writer.MainActivity"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContentView(R.layout.activity_main)

        val toolbar: Toolbar = findViewById(R.id.toolbar)
        setSupportActionBar(toolbar)
    }

    override fun onResume() {
        super.onResume()
        registerNFC()
    }

    override fun onPause() {
        super.onPause()
        unregisterNFC()
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)

        val nfcActions = arrayOf(
            NfcAdapter.ACTION_NDEF_DISCOVERED,
            NfcAdapter.ACTION_TECH_DISCOVERED,
            NfcAdapter.ACTION_TAG_DISCOVERED
        )

        if (intent.action in nfcActions) {
            onNewNFCIntent(intent)
        }
    }

    override fun onCreateOptionsMenu(menu: Menu): Boolean {
        menuInflater.inflate(R.menu.menu_main, menu)
        return true
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        return when (item.itemId) {
            R.id.action_settings -> {
                startActivity(Intent(this, SettingsActivity::class.java))
                true
            }
            else -> super.onOptionsItemSelected(item)
        }
    }



    private fun registerNFC() {
        val nfcAdapter = NfcAdapter.getDefaultAdapter(this)
        if (nfcAdapter != null) {
            val launchIntent = Intent(this, MainActivity::class.java).apply {
                addFlags(Intent.FLAG_ACTIVITY_SINGLE_TOP)
            }

            val pendingIntent = PendingIntent.getActivity(
                this,
                0,
                launchIntent,
                PendingIntent.FLAG_CANCEL_CURRENT or PendingIntent.FLAG_MUTABLE
            )

            nfcAdapter.enableForegroundDispatch(this, pendingIntent, null, null)
        }
    }

    /**
     * unregister NFC
     */
    private fun unregisterNFC() {
        val nfcAdapter = NfcAdapter.getDefaultAdapter(this)
        if (nfcAdapter != null) {
            nfcAdapter.disableForegroundDispatch(this)
        }
    }

    /**
     * Handle NFC intent
     */
    private fun onNewNFCIntent(intent: Intent) {
        Log.d(TAG, "NFC intent received: $intent")

        val tag = IntentCompat.getParcelableExtra(intent, NfcAdapter.EXTRA_TAG, Tag::class.java)
        if (tag == null) {
            Log.e(TAG, "No tag found in intent")
            return
        }

        writeNFC(tag)
    }


    /**
     * Write NFC tag
     */
    private fun writeNFC(tag: Tag) {
        Log.d(TAG, "Writing to tag: $tag")


        val iso = IsoDep.get(tag)
        Thread {
            try {
                iso.connect()

                val communicator = DnaCommunicator().apply {
                    setTransceiver { bytesToSend -> iso.transceive(bytesToSend) }
                    setLogger { info -> Log.d(TAG, "Communicator: $info") }
                    beginCommunication()
                }

                // Sync Key
                val keySet = getKeySet()
                keySet.synchronizeKeys(communicator)

                // Authenticate
                if (!AESEncryptionMode.authenticateEV2(communicator, 0, keySet.getKey(0).key)) {
                    Log.e(TAG, "Failed to authenticate")
                    runOnUiThread {
                        Toast.makeText(this, R.string.authentication_failed, Toast.LENGTH_SHORT).show()
                    }
                    return@Thread
                }

                Log.d(TAG, "Authenticated")

                val uid = GetCardUid.run(communicator)
                Log.d(TAG, "Card UID: ${ByteUtil.byteToHex(uid)}")

                val ndeffs = GetFileSettings.run(communicator, Ntag424.NDEF_FILE_NUMBER)

                val sdmSettings = SDMSettings().apply {
                    sdmMetaReadPerm = Permissions.ACCESS_EVERYONE
                    sdmFileReadPerm = Permissions.ACCESS_KEY2
                    sdmOptionUid = true
                    sdmOptionReadCounter = true
                }

                val master = NdefTemplateMaster().apply { usesLRP = false }
                val record = master.generateNdefTemplateFromUrlString("https://ntag424.aplulu.me/redeem?x={UID}{COUNTER}{MAC}", sdmSettings)

                WriteData.run(communicator, Ntag424.NDEF_FILE_NUMBER, record)

                ndeffs.apply {
                    readPerm = Permissions.ACCESS_EVERYONE
                    writePerm = Permissions.ACCESS_KEY0
                    readWritePerm = Permissions.ACCESS_KEY0
                    changePerm = Permissions.ACCESS_KEY0
                    this.sdmSettings = sdmSettings
                }

                ChangeFileSettings.run(communicator, Ntag424.NDEF_FILE_NUMBER, ndeffs)

                runOnUiThread {
                    Toast.makeText(this, R.string.tag_written_successfully, Toast.LENGTH_SHORT).show()
                    playSuccessSound()
                }
            } catch (e: Exception) {
                Log.e(TAG, "Error writing to tag", e)
                runOnUiThread {
                    Toast.makeText(this, R.string.tag_written_failed, Toast.LENGTH_SHORT).show()
                }
            } finally {
                iso.close()
                Log.d(TAG, "Connection closed")
            }
        }.start()
    }

    private fun getKeySet(): KeySet {
        val keySet = KeySet()
        keySet.usesLrp = false

        // Master Key
        val key0 = KeyInfo().apply {
            diversifyKeys = false
            key = Ntag424.FACTORY_KEY
        }
        keySet.setKey(Permissions.ACCESS_KEY0, key0)

        val key1 = KeyInfo().apply {
            diversifyKeys = false
            key = Ntag424.FACTORY_KEY
        }
        keySet.setKey(Permissions.ACCESS_KEY1, key1)

        val key2 = KeyInfo().apply {
            diversifyKeys = false
            key = Ntag424.FACTORY_KEY
        }
        keySet.setKey(Permissions.ACCESS_KEY2, key2)

        val key3 = KeyInfo().apply {
            diversifyKeys = false
            key = Ntag424.FACTORY_KEY
        }
        keySet.setKey(Permissions.ACCESS_KEY3, key3)

        val key4 = KeyInfo().apply {
            diversifyKeys = false
            key = Ntag424.FACTORY_KEY
        }
        keySet.setKey(Permissions.ACCESS_KEY4, key4)

        keySet.setMetaKey(Permissions.ACCESS_KEY2)
        keySet.setMacFileKey(Permissions.ACCESS_KEY3)

        return keySet
    }

    private fun playSuccessSound() {
        val mediaPlayer = MediaPlayer.create(this, R.raw.success_sound)
        mediaPlayer.start()
        mediaPlayer.setOnCompletionListener { mp -> mp.release() }
    }
}
