package com.xchat.stevenzack.langenius;

import android.content.DialogInterface;
import android.content.SharedPreferences;
import android.os.Build;
import android.support.v4.content.ContextCompat;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.text.InputType;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

public class SettingsActivity extends AppCompatActivity {
    private Button bt_filercvpath,bt_default_port;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_settings);
        final SharedPreferences sp_settings=getSharedPreferences(SettingsActivity.this.getString(R.string.sp_settings),MODE_PRIVATE);
        bt_filercvpath=(Button)findViewById(R.id.set_bt_frcvpath);
        bt_filercvpath.setText(this.getString(R.string.storagepath)+sp_settings.getString(this.getString(R.string.sp_sub_filercv_path),this.getString(R.string.default_filercvpath)));
        bt_filercvpath.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                AlertDialog.Builder builder = new AlertDialog.Builder(SettingsActivity.this);
                builder.setTitle(SettingsActivity.this.getString(R.string.str_setDefaultFileRecvPath));
                final EditText input = new EditText(SettingsActivity.this);
                input.setInputType(InputType.TYPE_CLASS_TEXT);
                input.setText(sp_settings.getString(SettingsActivity.this.getString(R.string.sp_sub_filercv_path),SettingsActivity.this.getString(R.string.default_filercvpath)));
                builder.setView(input);
                builder.setPositiveButton(SettingsActivity.this.getString(R.string.str_ok), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        String m_Text = input.getText().toString();
                        if (!m_Text.endsWith("/")){
                            m_Text=m_Text+"/";
                        }
                        sp_settings.edit().putString(SettingsActivity.this.getString(R.string.sp_sub_filercv_path),m_Text).commit();
                        bt_filercvpath.setText(SettingsActivity.this.getString(R.string.storagepath)+sp_settings.getString(SettingsActivity.this.getString(R.string.sp_sub_filercv_path),SettingsActivity.this.getString(R.string.default_filercvpath)));
                    }
                });
                builder.setNegativeButton(SettingsActivity.this.getString(R.string.str_cancel), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.cancel();
                    }
                });

                builder.show();
            }
        });
        bt_default_port=(Button)findViewById(R.id.set_bt_default_port);
        bt_default_port.setText(SettingsActivity.this.getString(R.string.str_default_port)+sp_settings.getString(this.getString(R.string.sp_sub_port),this.getString(R.string.default_port)));
        bt_default_port.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                AlertDialog.Builder builder = new AlertDialog.Builder(SettingsActivity.this);
                builder.setTitle(SettingsActivity.this.getString(R.string.str_setDefaultPort));
                final EditText input = new EditText(SettingsActivity.this);
                input.setInputType(InputType.TYPE_CLASS_NUMBER);
                input.setHint("1000~65535");
                builder.setView(input);
                builder.setPositiveButton(SettingsActivity.this.getString(R.string.str_ok), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        String m_Text = ":"+input.getText().toString();
                        sp_settings.edit().putString(SettingsActivity.this.getString(R.string.sp_sub_port),m_Text).commit();
                        bt_default_port.setText(SettingsActivity.this.getString(R.string.str_default_port)+sp_settings.getString(SettingsActivity.this.getString(R.string.sp_sub_port),SettingsActivity.this.getString(R.string.default_port)));
                    }
                });
                builder.setNegativeButton(SettingsActivity.this.getString(R.string.str_cancel), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.cancel();
                    }
                });

                builder.show();
            }
        });
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            getWindow().setNavigationBarColor(ContextCompat.getColor(this,R.color.colorPrimaryDark));
        }
    }
}
