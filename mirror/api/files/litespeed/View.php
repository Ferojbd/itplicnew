<?php

/** ******************************************
 * LiteSpeed Web Server Plugin for Plesk Panel
 *
 * @author    LiteSpeed Technologies, Inc. (https://www.litespeedtech.com)
 * @copyright 2013-2023
 * ******************************************* */

use Lsc\Wp\WPInstallStorage;
use Modules_Litespeed_PleskPluginException as PleskPluginException;
use Modules_Litespeed_Util as Util;
use Modules_Litespeed_View_Model_MissingTplViewModel as MissingTplViewModel;
use Modules_Litespeed_View_View as TplView;

class Modules_Litespeed_View
{

    /**
     * @var string[]
     */
    private $bufs = array();

    /**
     * @since 1.4
     * @var   string
     */
    private $do = 'main';

    /**
     * @deprecated 1.4  Going to be made private in the future.
     * @var        string[]
     */
    public $icons;

    /**
     * @var Zend_View
     */
    private $_parentview;

    /**
     *
     * @param Zend_View $view
     */
    public function __construct( $view )
    {
        $this->_parentview = $view;

        $this->icons = array(
            'm_logo_lsws'           => 'images/Logo_centered.svg',
            'm_server_version'      => 'icons/lsCurrentVersion.svg',
            'm_server_install'      => 'icons/install.png',
            'm_server_uninstall'    => 'icons/uninstall.svg',
            'm_server_definehome'   => 'icons/lsws-home.png',
            'm_server_php'          => 'icons/LSPHP.png',
            'm_cache_manage'        => 'icons/manageCacheInstallations.svg',
            'm_cache_mass_op'       => 'icons/massEnableDisableCache.svg',
            'm_cache_ver_manage'    => 'icons/lscwpCurrentVersion.svg',
            'm_dash_notifier'       => 'icons/wpNotifier.svg',
            'm_control_config'      => 'icons/lsConfiguration.svg',
            'm_control_restart'     => 'icons/restartLs.svg',
            'm_control_restart_php' => 'icons/restartDetachedPHP.svg',
            'm_license_check'       => 'icons/licenseStatus.svg',
            'm_license_change'      => 'icons/changeLicense.svg',
            'm_license_transfer'    => 'icons/transferLicense.svg',
            'm_switch_apache'       => 'icons/switchToApache.svg',
            'm_switch_lsws'         => 'icons/switchToLiteSpeed.svg',
            'm_switch_port'         => 'icons/changePortOffset.svg',
            'v_active'              => 'icons/ok.png',
            'ico_info'              => 'icons/server-info.png',
            'ico_error'             => 'icons/error.png',
        );
    }

    /**
     * Displays all HTML output stored in $this->bufs[].
     */
    public function dispatch()
    {
        echo implode("\n", $this->bufs);
    }

    public function returnAjax()
    {
        $this->dispatch();
        exit(0);
    }

    /**
     *
     * @since 1.4     Deprecated param $do.
     * @since 1.4.17  Removed deprecated parameter $do.
     */
    public function PageHeader()
    {
        if ( func_num_args() == 1 ) {
            trigger_error(
                'Parameter $do is no longer used and will be removed in a '
                    . 'future version.',
                E_USER_DEPRECATED
            );
        }

        $this->bufs[] = <<<EEN
<div id="lswsContent">
  <div id="lsws-container">
    <form name="lswsform">
      <input type="hidden" name="step" value="1" />
      <input type="hidden" name="do" value="$this->do" />
      <div class="form-box">
EEN;
    }

    public function PageFooter()
    {
        $this->bufs[] = "</div></form></div></div>";
    }

    /**
     *
     * @since 1.4
     *
     * @param string $do
     */
    public function setDoValue( $do )
    {
        $this->do = $do;
    }

    public function moduleError()
    {
        $this->PageHeader();

        $this->bufs[] =
            $this->screen_title(
                'Complete LiteSpeed Extension Installation'
            )
            . $this->error_panel_msg(
                  'Module is not installed properly',
                  'Please download and reinstall the extension from zip file.'
            )
        ;

        $this->PageFooter();
    }

    /**
     *
     * @param string[][] $list
     *
     * @return string
     */
    private function tool_list_block( array $list )
    {
        $buf = '<div class="tools-list">';

        foreach ( $list as $li ) {
            $buf .= '<div class="item" role="button"><a class="item-link" '
                . (($li['link']) ? "href=\"{$li['link']}\" " : '')
            ;

            if ( substr($li['link'], 0, 4) == 'http' ) {
                $buf .= 'target="_blank" rel="noopener noreferrer" ';
            }

            $buf .= 'title="'
                . (($li['info'] != '') ? $li['info'] : $li['name'])
                . '">'
            ;

            if ( $li['icon'] != '' ) {
                $buf .= "<img class=\"itemImageWrapper\" src=\"{$li['icon']} "
                    . 'alt="" />';
            }

            $buf .= "<span class=\"itemTextWrapper\">{$li['name']}</span>"
                . '</a></div>';
        }

        return "$buf</div>\n";
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function RestartLswsConfirm( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Confirm Operation... Restart LiteSpeed',
                true,
                $this->icons['m_control_restart']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        $msg = array();

        if ( $info['port_offset'] != 0 ) {
            $msg[] = "Apache port offset is {$info['port_offset']}.";
            $msg[] = 'LiteSpeed will be running in parallel with Apache. When '
                . 'you are ready to replace Apache with LiteSpeed, use the '
                . '<b>Switch to LiteSpeed</b> option.';
        }

        $goNext = 'Restart';

        if ( (int)$info['ap_pid'] > 0 && $info['port_offset'] == 0 ) {
            $msg[] = 'Apache port offset is 0. If you wish to use LiteSpeed as '
                . 'your main web server, please use the '
                . '<b>Switch to LiteSpeed</b> option.';
            $msg[] = 'If you need to run LiteSpeed in parallel with Apache, '
                . 'please use the <b>Change Port Offset</b> option.';
            $goNext = '';
        }

        if ( $goNext == 'Restart' ) {
            $msg[] = 'This will do a graceful restart of LiteSpeed Web Server.';
        }

        $this->bufs[] = $buf
            . $this->info_panel_msg(null, $msg)
            . $this->button_panel_back_next('Cancel', $goNext)
        ;

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function RestartLsws( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Restart LiteSpeed',
                true,
                $this->icons['m_control_restart']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        if ( (int)$info['ls_pid'] > 0 ) {
            $buf .= $this->info_panel_msg(
                'LiteSpeed restarted successfully',
                $info['output']
            );
        }
        else {
            $buf .= $this->error_panel_msg(
                'LiteSpeed is not running! Please check the web server log '
                    . 'file for errors.',
                $info['output']
            );
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function Switch2LswsConfirm( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Confirm Operation... Switch to LiteSpeed',
                true,
                $this->icons['m_switch_lsws']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        if ( $info['port_offset'] != 0 ) {
            $buf .= $this->warning_msg(
                "Apache port offset is {$info['port_offset']}. This action "
                . 'will change port offset to 0.'
            );
        }

        $msg = array();

        if ( (int)$info['ap_pid'] > 0 ) {
            $msg[] = 'This action will stop Apache and restart LiteSpeed as '
                . 'the main web server. It may take a little while for Apache '
                . 'to stop completely.';
        }

        $msg[] = 'This will restart '
            . '<strong>LiteSpeed as main web server</strong>!';

        $this->bufs[] = $buf
            . $this->info_panel_msg(null, $msg)
            . $this->button_panel_back_next('Cancel', 'Switch to LiteSpeed')
        ;

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function Switch2Lsws( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Switch To LiteSpeed',
                true,
                $this->icons['m_switch_lsws']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        $out = $info['output'];

        if ( $info['port_offset'] != 0 ) {
            $out[] = 'Failed to set Apache port offset to 0. Please check '
                . 'config file permissions.';
        }
        else {
            $out[] = 'Apache port offset has been set to 0.';
        }

        if ( (int)$info['ls_pid'] > 0 ) {
            $buf .= $this->info_panel_msg(
                'Switched to LiteSpeed successfully',
                $out
            );
        }
        else {
            $buf .= $this->error_panel_msg(
                'Failed to bring up LiteSpeed',
                $out
            );
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function Switch2ApacheConfirm( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Confirm Operation... Switch to Apache',
                true,
                $this->icons['m_switch_apache']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        if ( isset($info['stop_msg']) ) {
            $msg       = $info['stop_msg'];
            $backTitle = 'OK';
            $nextTitle = '';
        }
        else {
            $msg = array();

            if ( (int)$info['ls_pid'] > 0 ) {
                $msg[] = 'This action will stop LiteSpeed and restart Apache '
                    . 'as the main web server. It may take a little while for '
                    . 'LiteSpeed to stop completely.'
                ;
            }

            $msg[] = 'This will restart '
                . '<strong>Apache as main web server</strong>!';

            $backTitle = 'Cancel';
            $nextTitle = 'Switch to Apache';
        }

        $this->bufs[] = $buf
            . $this->info_panel_msg(null, $msg)
            . $this->button_panel_back_next($backTitle, $nextTitle)
        ;

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function Switch2Apache( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Switch To Apache',
                true,
                $this->icons['m_switch_apache']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        if ( isset($info['stop_msg']) ) {
            $buf .= $this->info_panel_msg(null, $info['stop_msg']);
        }
        elseif ( $info['return'] != 0 ) {
            $buf .= $this->info_panel_msg(null, $info['output'])
                . $this->error_panel_msg(
                    'Failed to switch to Apache',
                    array(
                        'Failed to switch to Apache!',
                        'This may be due to a configuration error. To manually '
                            . 'check this problem, please ssh to your server.',
                        'Use the following steps to manually switch to Apache:',
                        'Stop LiteSpeed if lshttpd still running: '
                            . '<code>pkill -9 lshttpd </code>',
                        'Restore Apache httpd if /usr/sbin/httpd_ls_bak '
                            . 'exists: <code>'
                            . 'mv -f /usr/sbin/httpd_ls_bak /usr/sbin/httpd'
                            . '</code>',
                        'Run the Apache restart command manually: '
                            . '<code>service httpd restart</code> and check '
                            . 'for errors.'
                    )
                )
            ;
        }
        else {
            $info['output'][] = 'If you\'d like to re-enable Nginx reverse '
                . 'proxy server while running Apache, you can do so from '
                . '"Server Management > Tools & Settings > Services '
                . 'Management".';

            $buf .= $this->info_panel_msg(
                'Switched to Apache successfully',
                $info['output']
            );
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function ChangePortOffsetConfirm( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Confirm Operation... Change LiteSpeed Port Offset',
                true,
                $this->icons['m_switch_port']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
            . '<div class="indent-box"><div class="hint"><p>Port offset allows '
            . 'you to run Apache and LiteSpeed in parallel by running '
            . 'LiteSpeed on a separate port.</p><p>For example, if Apache is '
            . 'running on port 80 and the LiteSpeed port offset is 2000, then '
            . 'you will be able to access LiteSpeed-powered web pages on port '
            . '2080.</p></div></div>'
        ;

        if ( $info['port_offset'] == 0 && (int)$info['ap_pid'] == 0 ) {
            $buf .= $this->warning_msg(
                array(
                    'Apache is currently not running. We suggest your first '
                        . '<strong>switch to Apache </strong> to avoid server '
                        . 'downtime.'
                )
            );
        }

        $this->bufs[] = $buf
            . $this->section_title('Change Port Offset')
            . $this->form_row(
                'Set new port offset',
                'text',
                $this->input_text('port_offset', null),
                ((isset($info['error'])) ? $info['error'] : null),
                "Current Port Offset is {$info['port_offset']}."
            )
            . $this->button_panel_back_next('Cancel', 'Change')
        ;

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     */
    public function ChangePortOffset( array $info )
    {
        $this->PageHeader();

        $buf = $this->screen_title(
            'Change LiteSpeed Port Offset',
            true,
            $this->icons['m_switch_port']
        );

        if ( $info['return'] != 0 ) {
            $buf .= $this->error_panel_msg(
                'Failed to Change Port Offset',
                $info['output']
            );
        }
        else {
            $buf .= $this->info_panel_msg(
                'Saved New Port Offset',
                "Port offset has been changed to {$info['new_port_offset']}."
            );
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     */
    public function CheckLicense( array $info )
    {
        $this->PageHeader();

        $buf = $this->screen_title(
            "Current License Status <small>(Serial: {$info['serial']})</small>",
            true,
            $this->icons['m_license_check']
        );

        $lic_type = empty($info['lic_type']) ? '' : "{$info['lic_type']} - ";

        if ( $info['return'] != 0 ) {
            $buf .= $this->error_panel_msg(
                "{$lic_type}Error when checking license status",
                $info['output']
            );
        }
        else {
            $buf .= $this->info_panel_msg(
                "{$lic_type}Check against license server",
                $info['output']
            );
        }

        if ( $info['lictype'] == 'trial' ) {
            $buf .= $this->info_msg(
                'Note: For trial licenses, the expiration date above is based '
                    . 'on the most recent trial license you have downloaded. '
                    . 'All trial licenses are valid for 15 days from the day '
                    . 'you apply. Each IP address, though, may only use trial '
                    . 'licenses for 30 days from the date of the first '
                    . 'application. The expiration date above does not reflect '
                    . 'how much longer your IP may use trial licenses. If you '
                    . 'are on your second or third trial license, your actual '
                    . 'trial period may end earlier than the above date.'
            );
        }
        elseif ( $info['lictype'] == 'migrated' ) {
            $buf .= $this->warning_msg(
                'Note: You have started the license migration process. You can '
                    . 'now use the same serial number to register on a new '
                    . 'machine. If you decide you want to continue using the '
                    . 'license on this machine instead, you must re-register '
                    . 'the license here. Use the Change License function with '
                    . 'the serial number to re-register.'
            );
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function UninstallLswsPrepare( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Confirm Operation... Uninstall LiteSpeed Web Server',
                true,
                $this->icons['m_server_uninstall']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        if ( isset($info['stop_msg']) ) {
            $this->bufs[] = $buf
                . $this->error_msg($info['stop_msg'])
                . $this->button_panel_back_next('OK', '', 'config_lsws')
            ;
        }
        else {
            $msg = array();

            if ( (int)$info['ls_pid'] > 0 ) {
                $msg[] = 'LiteSpeed is currently running on port offset '
                    . "{$info['port_offset']} and will be stopped first.";
            }

            $msg[] = 'All subdirectories created under ' . Util::getLswsHome()
                . ' during installation will be removed! The conf/ and logs/ '
                . 'subdirectories can be preserved using the check boxes '
                . 'below.';
            $msg[] = 'If you want to preserve any files under other '
                . 'subdirectories created by the installation script, please '
                . 'manually back them up before proceeding.';

            $this->bufs[] = $buf
                . $this->section_title('Uninstall Options')
                . $this->warning_msg($msg)
                . $this->form_row(
                    'Keep conf/ directory',
                    '',
                    $this->input_checkbox('keep_conf', '1', true),
                    null,
                    null,
                    true
                )
                . $this->form_row(
                    'Keep logs/ directory',
                    '',
                    $this->input_checkbox('keep_log', '1', true),
                    null,
                    null,
                    true
                )
                . $this->button_panel_back_next(
                    'Cancel',
                    'Uninstall',
                    'config_lsws'
                )
            ;
        }

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function UninstallLsws( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Uninstall LiteSpeed Web Server',
                true,
                $this->icons['m_server_uninstall']
            )
            . $this->show_running_status($info, false)
        ;

        if ( $info['return'] != 0 ) {
            $buf .= (($info['spool_warning']) ? $this->getSpoolWarning() : '')
                . $this->error_panel_msg(
                    'Error when uninstalling LiteSpeed',
                    $info['output']
                )
            ;
        }
        else {
            $buf .= $this->info_panel_msg(
                'Uninstalled successfully',
                $info['output']
            );
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @return string
     */
    private function show_choose_license( array $info )
    {
        $licAgreeErr = $installTypeErr = $serialNoErr = null;

        if ( isset($info['error']['license_agree']) ) {
            $licAgreeErr = $info['error']['license_agree'];
        }

        if ( isset($info['error']['install_type']) ) {
            $installTypeErr = $info['error']['install_type'];
        }

        if ( isset($info['error']['serial_no']) ) {
            $serialNoErr = $info['error']['serial_no'];
        }

        /** @noinspection HtmlUnknownTarget */
        return '<div>'
            . '<iframe src="LICENSE.html" width="650" height="400"></iframe>'
            . '</div>'
            . $this->form_row(
                'I agree',
                '',
                $this->input_checkbox(
                    'license_agree',
                    'agree',
                    ($info['license_agree'] == 'agree')
                ),
                $licAgreeErr,
                null,
                true
            )
            . $this->section_title('Choose a License Type')
            . $this->form_row(
                'Use an Enterprise license',
                'radio',
                $this->installTypeInputRadio(
                    'prod',
                    ($info['install_type'] == 'prod')
                ),
                $installTypeErr,
                null,
                true
            )
            . $this->form_row(
                'Input serial number:',
                'text',
                $this->input_text('serial_no', $info['serial_no'], 1),
                $serialNoErr,
                array(
                    'Your serial number is sent via email when you purchase a '
                        . 'LiteSpeed Web Server license. You can also copy it '
                        . 'from your service details in our client area '
                        . '(store.litespeedtech.com).'
                )
            )
            . $this->warning_msg(
                'If your license is currently running on another server, you '
                    . 'will need to transfer the license (using the Transfer '
                    . 'License function) before registering it on this server.'
            )
            . $this->form_row(
                'Request a trial license',
                'radio',
                $this->installTypeInputRadio(
                    'trial',
                    ($info['install_type'] == 'trial')
                ),
                $installTypeErr,
                null,
                true
            )
            . $this->form_row(
                '',
                '',
                '',
                null,
                array(
                    'This will retrieve a trial license from the LiteSpeed '
                        . 'license server.',
                    'Each trial license is valid for 15 days from the date you '
                        . 'apply.',
                    'Each IP address can only use trial licenses for 30 days '
                        . 'from the date of your first application.',
                    'If you need to extend your trial period, please contact '
                        . 'the sales department at litespeedtech.com.'
                )
            )
        ;
    }

    /**
     *
     * @param array $info
     */
    public function ChangeLicensePrepare( array $info )
    {
        $this->PageHeader();

        $this->bufs[] =
            $this->screen_title(
                'Changing Software License for LiteSpeed Web Server',
                true,
                $this->icons['m_license_change']
            )
            . $this->show_choose_license($info)
            . $this->button_panel_back_next('Cancel', 'Switch')
        ;

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function ChangeLicense( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Changing Software License for LiteSpeed Web Server',
                true,
                $this->icons['m_license_change']
            )
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        if ( $info['return'] != 0 ) {
            $buf .= $this->error_panel_msg(
                'Error when activating new license',
                $info['output']
            );
        }
        else {
            $buf .= $this->info_panel_msg(
                'New license activated successfully',
                $info['output']
            );
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     */
    public function TransferLicenseConfirm( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'LiteSpeed Web Server License Migration Confirm',
                true,
                $this->icons['m_license_transfer']
            )
            . $this->info_panel_msg(
                null,
                'You can transfer your license from one server to another. '
                    . 'This migration process will allow you to continue to '
                    . 'use your current server for 3 days while you migrate to '
                    . 'your new server. If, after 3 days, you still need more '
                    . 'time to use LiteSpeed on this server, please download a '
                    . '15-day trial license. (You will need to contact '
                    . 'LiteSpeed Technologies to reset your trial record if '
                    . 'this server has previously used trial licenses.)'
            )
            . $this->info_panel_msg(
                'Current license status',
                $info['licstatus_output']
            )
        ;

        if ( $info['error'] != '' ) {
            $this->bufs[] = $buf
                . $this->error_msg($info['error'])
                . $this->button_panel_back_next('OK')
            ;
        }
        else {
            $this->bufs[] = $buf
                . $this->warning_msg(
                    'Click Transfer if you are ready to go ahead and transfer '
                        . 'your current license. You can continue using this '
                        . 'server for up to 3 days.'
                )
                . $this->button_panel_back_next('Cancel', 'Transfer')
            ;
        }

        $this->PageFooter();
    }

    /**
     *
     * @param array $info
     */
    public function TransferLicense( array $info )
    {
        $this->PageHeader();

        $buf = $this->screen_title(
            'LiteSpeed Web Server License Migration',
            true,
            $this->icons['m_license_transfer']
        );

        if ( $info['return'] == 0 ) {
            $buf .= $this->info_panel_msg(
                'Successfully migrated your license',
                $info['output']
            );
        }
        else {
            $buf .= $this->error_panel_msg(
                'Failed to migrate your license',
                $info['output']
            );
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @since 1.4.17  Inlined unused parameter $scrollable.
     *
     * @param string|string[] $msg
     * @param string          $subtype
     * @param string          $title
     *
     * @return string
     */
    private function div_msg_box( $msg, $subtype, $title )
    {
        if ( empty($msg) ) {
            return '';
        }

        $div = '<div class="msg-box'
            . (($subtype != '') ? " $subtype" : '')
            . '"><ul><li>'
        ;

        if ( $title != '' ) {
            $div .= "<div class=\"title\">$title</div></li><li>";
        }

        return $div
            . ((is_array($msg)) ? implode('</li><li>', $msg) : $msg)
            . '</li></ul></div>'
        ;
    }

    /**
     *
     * @since 1.4.17  Inlined unused parameters $title and $scrollable.
     *
     * @param string|string[] $msg
     *
     * @return string
     */
    private function status_msg( $msg )
    {
        return $this->div_msg_box($msg, 'msg-status', '');
    }

    /**
     *
     * @since 1.4.17  Inlined unused parameters $title and $scrollable.
     *
     * @param string|string[] $msg
     *
     * @return string
     */
    private function info_msg( $msg )
    {
        return $this->div_msg_box($msg, 'msg-info', '');
    }

    /**
     *
     * @since 1.4.17  Inlined unused parameter $scrollable.
     *
     * @param string|string[] $msg
     * @param string          $title
     *
     * @return string
     */
    private function error_msg( $msg, $title = '' )
    {
        return $this->div_msg_box($msg, 'msg-error', $title);
    }

    /**
     *
     * @since 1.4.17  Inlined unused parameters $title and $scrollable.
     *
     * @param string|string[] $msg
     *
     * @return string
     */
    private function warning_msg( $msg )
    {
        return $this->div_msg_box($msg, 'msg-warn', '');
    }

    /**
     *
     * @param string|null     $title
     * @param string|string[] $msg
     *
     * @return string
     */
    private function info_panel_msg( $title, $msg )
    {
        return $this->div_msg_panel($title, $msg, $this->icons['ico_info']);
    }

    /**
     *
     * @param string|null     $title
     * @param string|string[] $msg
     *
     * @return string
     */
    private function error_panel_msg( $title, $msg )
    {
        return $this->div_msg_panel($title, $msg, $this->icons['ico_error']);
    }

    /**
     *
     * @param string|null     $title
     * @param string|string[] $msg
     * @param string          $icon
     *
     * @return string
     */
    private function div_msg_panel( $title, $msg, $icon )
    {
        $box = '<div class="p-box"><div class="p-box-content">';

        if ( $title != null ) {
            $box .= '<div class="title"><div class="title-area"><h4>'
                . "<img src=\"$icon\"  alt=\"\"/> $title</h4><p></p>"
                . '</div></div>';
        }

        return $box
            . '<div class="content"><div class="content-area"><p>'
            . ((is_array($msg)) ? implode('</p><p>', $msg) : $msg)
            . "</p></div></div></div></div>\n"
        ;
    }

    /**
     *
     * @since 1.4.3  Added optional parameter $getCriticalAlertMsg.
     *
     * @param array $info
     * @param bool  $getCriticalAlertMsg
     *
     * @return string
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     Util::getCriticalAlertMsg() call.
     */
    private function show_running_status(
        array $info,
              $getCriticalAlertMsg = true )
    {
        $ret = '';

        if ( $getCriticalAlertMsg ) {
            $ret .=
                $this->error_msg(Util::getCriticalAlertMsg(), 'Critical Alert');
        }

        $lsPid = (int)$info['ls_pid'];
        $apPid = (int)$info['ap_pid'];

        if ( $lsPid > 0 ) {
            $msg = "LiteSpeed is running (PID = $lsPid";

            if ( isset($info['port_offset']) ) {
                $msg .= ", Apache_Port_Offset = {$info['port_offset']}";
            }

            $msg .= ').';
        }
        else {
            $msg = 'LiteSpeed is not running.';
        }

        if ( $apPid > 0 ) {
            $msg .= " Apache is running (PID = $apPid).";
        }
        else {
            $msg .= ' Apache is not running.';
        }

        $ret .= $this->status_msg($msg);

        $showNginxReverseProxyMsg = (
            $lsPid > 0
            &&
            isset($info['nginx_proxy_running'])
            &&
            $info['nginx_proxy_running'] != ''
        );

        if ( $showNginxReverseProxyMsg ) {
            $ret .= '<div id="nginxNotice" class="msg-box msg-error"><ul><li>'
                . 'Nginx reverse proxy server is currently running and must be '
                . 'stopped.<span id="stopNginxBtn"> '
                . '<a role="button" onclick="stopNginxProxy();">Stop Now</a>'
                . '</span></li></ul></div>'
            ;
        }

        return $ret;
    }

    /**
     * Gets spool warning message HTML.
     *
     * Note: this can be moved to a Logger msg later on.
     *
     * @return string
     */
    private function getSpoolWarning()
    {
        return $this->warning_msg(
            'Both LiteSpeed and Apache are running with '
                . 'Apache_Port_Offset = 0. This can cause unintended server '
                . "behavior. Please do one of the following: "
                . '<a href="?do=switch_lsws" title="Use LiteSpeed as main web '
                . 'server. This will update rc scripts.">'
                . 'Switch to LiteSpeed</a> '
                . '| <a href="?do=switch_apache" title="Use Apache as main web '
                . 'server. This will update rc scripts.">'
                . 'Switch to Apache</a> '
                . '| <a href="?do=change_port_offset" title="This allows '
                . 'LiteSpeed and Apache to run in parallel.">'
                . 'Change Port Offset</a> to a non-zero value.'
        );
    }

    /**
     *
     * @param int $dataFileError
     *
     * @return string
     */
    private function checkDataFile( $dataFileError )
    {
        switch ( $dataFileError ) {

            case WPInstallStorage::ERR_NOT_EXIST:
                $msg = '<p><b>Deploy LiteSpeed Cache plugin for WordPress</b>'
                    . '</p><p>If you have never deployed the LiteSpeed Cache '
                    . 'for WordPress plugin across discovered WordPress sites, '
                    . 'visit the <a href="?do=lscwp_manage">Manage Cache '
                    . 'Installations</a> page and perform a Scan to get '
                    . 'started. Sever-wide deployment can lead to significant '
                    . 'performance improvements for said sites, as well as an '
                    . 'overall reduction in server load.</p><p>If you\'ve been '
                    . 'here before, your data file may have been removed due '
                    . 'to a necessary data file update. Please perform a new '
                    . 'scan to rebuild the file. We apologize for any '
                    . 'inconvenience.</p>'
                ;
                break;

            case WPInstallStorage::ERR_VERSION_LOW:
                $msg = 'Cache Management scan data file format has been '
                    . 'changed for this version. Please perform a '
                    . '<a href="?do=lscwp_manage">Re-scan</a>'
                    . ' before attempting any Cache Management operations.'
                ;
                break;

            default:
                return '';
        }

        return $this->warning_msg($msg);
    }

    /**
     *
     * @param string $title
     * @param bool   $uplinkSelf
     * @param string $icon
     *
     * @return string
     */
    private function screen_title($title, $uplinkSelf = true, $icon = '' )
    {
        if ( $uplinkSelf ) {
            /** @noinspection PhpUndefinedFieldInspection */
            $this->_parentview->uplevelLink = pm_Context::getBaseUrl();
        }

        return '<div id="heading"><h2>'
            . (($icon != '') ? "<span><img src=\"$icon\" alt> </span>" : '')
            . "$title</h2></div>\n"
        ;
    }

    /**
     *
     * @param string $title
     *
     * @return string
     */
    private function section_title( $title )
    {
        return '<div><fieldset class="border-bottom">'
            . "<legend class=\"legend-spacing\">$title</legend></fieldset>"
            . "</div>"
        ;
    }

    /**
     *
     * @param string      $name
     * @param string|null $value
     * @param int         $size_class
     *
     * @return string
     */
    private function input_text( $name, $value, $size_class = 0 )
    {
        /**
         * size 0: default, size 1: f-middle-size, size 2: long
         */
        switch ( $size_class ) {

            case 1:
                $iclass = 'f-middle-size input-text';
                break;

            case 2:
                $iclass = '" size="90';
                break;

            default:
                $iclass = 'input-text';
        }

        return '<input type="text" '
            . "class=\"$iclass\" name=\"$name\" value=\"$value\"/>"
        ;
    }

    /**
     *
     * @param string $name
     * @param string $value
     *
     * @return string
     */
    private function input_password( $name, $value )
    {
        return '<input type="password" '
            . "name=\"$name\" value=\"$value\" class=\"ls-pass\"/>"
        ;
    }

    /**
     *
     * @param string $name
     * @param string $value
     * @param bool   $isChecked
     *
     * @return string
     */
    private function input_checkbox( $name, $value, $isChecked )
    {
        return '<input type="checkbox" class="checkbox" '
            . "name=\"$name\" value=\"$value\""
            . (($isChecked) ? ' checked="checked"' : '')
            . ' /><label class="ls-checkbox"></label>'
        ;
    }

    /**
     *
     * @since 1.4.17  Renamed from input_radio().
     * @since 1.4.17  Inlined unused parameter $name.
     *
     * @param string $value
     * @param bool   $isChecked
     *
     * @return string
     */
    private function installTypeInputRadio( $value, $isChecked)
    {
        return "<input type=\"radio\" name=\"install_type\" value=\"$value\""
            . (($isChecked) ? ' checked="checked"' : '')
            . ' />'
        ;
    }

    /**
     *
     * @param string               $label
     * @param string               $type
     * @param string               $field
     * @param string|null          $err
     * @param string|string[]|null $hints
     * @param bool                 $is_single
     *
     * @return string
     */
    private function form_row(
        $label,
        $type,
        $field,
        $err,
        $hints = null,
        $is_single = false )
    {
        switch($type) {

            case 'text':
            case 'pass':
                $labelClass = 'ls-text-indent';
                break;

            case 'checkbox':
                $labelClass = 'ls-checkbox';
                break;

            default:
                $labelClass = '';
        }

        $divClass = 'ls-form-row form-row';
        $errSpan  = '';
        $hintSpan = '';

        if ( $err != null ) {
            $divClass .= ' error';
            $errSpan   = "<span class=\"error-hint\">$err</span>";
        }

        if ( $hints != null ) {

            if ( is_array($hints) ) {
                $hintSpan =
                    '<span class="hint">' . implode('<br>', $hints) . '</span>';
            }
            else {
                $hintSpan = "<span class=\"hint\">$hints</span>";
            }
        }

        $div = "<div class=\"$divClass\">";

        if ( $is_single ) {
            $div .= "<div class=\"single-row\">$field<label"
                . ((!empty($labelClass)) ? " class=\"$labelClass\"" : '')
                . "><span class=\"ls-input-label-space\">$label</span></label>"
            ;
        }
        else {
            $div .= '<div class="field-name"><label'
                . ((!empty($labelClass)) ? " class=\"$labelClass\"" : '')
                . "><span class=\"ls-input-label-space\">$label</span></label>"
                . "</div><div class=\"field-value\">$field"
            ;
        }

        return $div . $errSpan . "$hintSpan</div></div>\n";
    }

    /**
     *
     * @since 1.4.17  Removed unused parameter $extra_class.
     *
     * @param string $back_title
     * @param string $next_title
     * @param string $back_do
     *
     * @return string
     */
    private function button_panel_back_next(
        $back_title,
        $next_title = '',
        $back_do = 'main' )
    {
        /** @noinspection JSUnresolvedReference */
        $buf = '<div class="btns-box">'
            . '<button class="input-button btn" '
            . "onclick=\"lswsform.do.value='$back_do';lswsform.submit();\">"
            . "$back_title</button>"
        ;

        if ( $next_title != '' ) {
            $buf .= '<button class="input-button btn action" type="submit">'
                . "$next_title</button>"
            ;
        }

        return $buf . '</div>';
    }

    /**
     * Creates a View object with the passed ViewModel and inserts the prepared
     * HTML output into $this->bufs[]. Actual HTML output occurs in
     * $this->dispatch().
     *
     * @param object $viewModel  ViewModel object.
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->displayMissingTpl() call.
     */
    public function prepareView( $viewModel )
    {
        $this->PageHeader();

        $tplView = new TplView($viewModel);

        try
        {
            ob_start();

            $tplView->display();

            $this->bufs[] = ob_get_clean();
        }
        catch ( PleskPluginException $e )
        {
            if ( ob_get_level() == 2 ) {
                ob_clean();
            }

            ob_start();

            $this->displayMissingTpl($e->getMessage());

            $this->bufs[] = ob_get_clean();
        }

        $this->PageFooter();
    }

    /**
     *
     * @since 1.3
     *
     * @param string $msg
     *
     * @throws PleskPluginException  Thrown indirectly by $view->display() call.
     */
    private function displayMissingTpl( $msg )
    {
        $view = new TplView(new MissingTplViewModel($msg));
        $view->display();
    }

    /**
     * Functions below this comment have been deprecated.
     */

    /**
     *
     * @deprecated 1.4
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function MainMenu( array $info )
    {
        $this->PageHeader();

        $this->screen_title("LiteSpeed Extension", false);

        $buf = '<div id="heading"><div class="center"><img class="header-logo" '
            . 'alt="LiteSpeed Web Server" '
            . "src=\"{$this->icons['m_logo_lsws']}\" "
            . 'onclick="window.open(\'https://www.litespeedtech.com\')" >'
            . '</div></div>'
            . $this->show_running_status($info)
            . (($info['spool_warning']) ? $this->getSpoolWarning() : '')
        ;

        if ( isset($info['data_file_error']) ) {
            $buf .= $this->checkDataFile($info['data_file_error']);
        }

        $buf .= '<div id="main" class="clearfix">';

        if ( $info['is_installed'] ) {
            $buf .= $this->section_title('Manage LiteSpeed Web Server');

            $li_version = array(
                'icon' => $this->icons['m_server_version'],
                'link' => '?do=versionManager',
                'name' => "Current Version: {$info['lsws_version']}",
                'info' =>
                    'Version Management: upgrade/downgrade/force reinstall.'
            );

            if ( !empty($info['lsws_build']) ) {
                $li_version['name'] .= " (build {$info['lsws_build']})";
            }

            $li_version['name'] .= '<div class="release-alert small red">';

            if ( !empty($info['new_build']) ) {
                $li_version['name'] .=
                    "<span>Latest Build: {$info['new_build']}</span><br />";
            }

            if ( $info['new_version'] != ''
                && $info['new_version'] != $info['lsws_version'] ) {

                $li_version['name'] .=
                    "<span>Latest Release: {$info['new_version']}</span>";
            }

            $li_version['name'] .= '</div>';

            $buf .= $this->tool_list_block(
                array(
                    $li_version,
                    array(
                        'icon' => $this->icons['m_control_config'],
                        'link' => '?do=config_lsws',
                        'name' => 'LiteSpeed Configuration',
                        'info' => 'Configure LiteSpeed settings.'
                    ),
                    array(
                        'icon' => $this->icons['m_control_restart'],
                        'link' => '?do=restart_lsws',
                        'name' => 'Restart LiteSpeed',
                        'info' => 'Gracefully restart LiteSpeed.'
                    ),
                    array(
                        'icon' => $this->icons['m_control_restart_php'],
                        'link' => '?do=restart_detached_php',
                        'name' => 'Restart Detached PHP Processes',
                        'info' => ''
                    )
                )
            );

            $cache_check = $manage_link = $massEnableDisable_link
                = $verManager_link = '';

            switch ( $info['has_cache'] ) {

                case Util::LSCACHE_STATUS_UNKNOWN:
                    $cache_check =
                        '(Please start LiteSpeed to access these features)';
                    break;

                case Util::LSCACHE_STATUS_MISSING:
                    $cache_check = '(This feature requires '
                        . '<a href="https://docs.litespeedtech.com/'
                        . 'licenses/how-to/#add-cache-to-an-existing-license" '
                        . 'target="_blank" '
                        . 'rel="noopener noreferrer">LSCache</a>)'
                    ;
                    break;

                case Util::LSCACHE_STATUS_DETECTED:
                    $manage_link            = '?do=lscwp_manage';
                    $massEnableDisable_link = '?do=lscwp_mass_enable_disable';
                    $verManager_link        = '?do=lscwpVersionManager';
                    break;

                //no default
            }

            if ( $info['has_cache'] != Util::LSCACHE_STATUS_NOT_SUPPORTED ) {
                $buf .= $this->section_title(
                    "LiteSpeed Cache For WordPress Management $cache_check"
                );

                $ver_mgr = array(
                    'icon' => $this->icons['m_cache_ver_manage'],
                    'link' => $verManager_link,
                    'info' => 'Change active cache plugin version or force a '
                        . 'version change for existing installations.'
                );

                if ( !$info['lscwp_curr_ver'] ) {
                    $ver_mgr['name'] = 'LSCWP Version Manager';
                }
                else {
                    $ver_mgr['name'] = 'Current Version: '
                        . htmlspecialchars($info['lscwp_curr_ver']);

                    $newerLscwpVerAvailable = (
                        $info['lscwp_latest']
                        &&
                        $info['lscwp_latest'] != $info['lscwp_curr_ver']
                    );

                    if ( $newerLscwpVerAvailable ) {
                        $ver_mgr['name'] .= '<br /><span class="small red">'
                            . 'Latest Release: '
                            . htmlspecialchars($info['lscwp_latest'])
                            . '</span>'
                        ;
                    }
                }

                $buf .= $this->tool_list_block(
                    array(
                        array(
                            'icon' => $this->icons['m_cache_manage'],
                            'link' => $manage_link,
                            'name' => 'Manage Cache Installations',
                            'info' => 'Enable/Disable cache or set a flag for '
                                . 'known WordPress installations.'
                        ),
                        array(
                            'icon' => $this->icons['m_cache_mass_op'],
                            'link' => $massEnableDisable_link,
                            'name' => 'Mass Enable/Disable Cache',
                            'info' => 'Enable/Disable cache for all '
                                . 'non-flagged WordPress installations'
                        ),
                        $ver_mgr,
                        array(
                            'icon' => $this->icons['m_dash_notifier'],
                            'link' => '?do=dash_notifier',
                            'name' => 'WordPress Dash Notifier',
                            'info' => 'Recommend a plugin or broadcast a '
                                . 'message to all discovered WordPress '
                                . 'installations.'
                        )
                    )
                );
            }

            $buf .= $this->section_title('License Management');

            if ( $info['serial'] == 'TRIAL' ) {
                $licenseInfo = '15-Day Trial License';
            }
            else {
                $licenseInfo = $info['serial'];
            }

            $list = array(
                array(
                    'icon' => $this->icons['m_license_check'],
                    'link' => '?do=check_current_license',
                    'name' => 'License Status <br/>'
                        . "<span class=\"small cornflower-blue\">$licenseInfo"
                        . '</span>',
                    'info' => 'Check/Refresh current license.'
                ),
                array(
                    'icon' => $this->icons['m_license_change'],
                    'link' => '?do=change_license',
                    'name' => 'Change License',
                    'info' => 'Switch to another license.'
                )
            );

            if ( $info['serial'] != 'TRIAL' ) {
                $list[] = array(
                    'icon' => $this->icons['m_license_transfer'],
                    'link' => '?do=transfer_license',
                    'name' => 'Transfer License',
                    'info' => 'Start license migration. Frees license for '
                        . 'registration on another server while leaving '
                        . 'license active on the current server for a limited '
                        . 'time.'
                );
            }

            $buf .= $this->tool_list_block($list)
                . $this->section_title('Switch between Apache and LiteSpeed')
                . $this->tool_list_block(
                    array(
                        array(
                            'icon' => $this->icons['m_switch_apache'],
                            'link' => '?do=switch_apache',
                            'name' => 'Switch to Apache',
                            'info' => 'Use Apache as main web server. This '
                                . 'will update rc scripts.'
                        ),
                        array(
                            'icon' => $this->icons['m_switch_lsws'],
                            'link' => '?do=switch_lsws',
                            'name' => 'Switch to LiteSpeed',
                            'info' => 'Use LiteSpeed as main web server. This '
                                . 'will update rc scripts.'
                        ),
                        array(
                            'icon' => $this->icons['m_switch_port'],
                            'link' => '?do=change_port_offset',
                            'name' => 'Change Port Offset',
                            'info' => 'LiteSpeed port offset is '
                                . "{$info['port_offset']}. This allows "
                                . 'LiteSpeed and Apache to run in parallel.'
                        )
                    )
                )
            ;
        }
        else {
            $buf .= $this->section_title('Install LiteSpeed Web Server')
                . $this->tool_list_block(
                    array(
                        array(
                            'icon' => $this->icons['m_server_install'],
                            'link' => '?do=install',
                            'name' => 'Install LiteSpeed Web Server',
                            'info' => 'Download and install the latest stable '
                                . 'release.'
                        ),
                        array(
                            'icon' => $this->icons['m_server_definehome'],
                            'link' => '?do=define_home',
                            'name' => 'Define LSWS_HOME',
                            'info' => 'If you installed LiteSpeed Web Server '
                                . 'before installing this extension, please '
                                . 'specify your LSWS_HOME directory before '
                                . 'using the extension.'
                        )
                    )
                )
            ;
        }

        $this->bufs[] = $buf
            . '<p style="margin-top:30px;color:#a0a0a0;text-align:right;'
            . 'font-size:11px">This extension is developed by LiteSpeed '
            . 'Technologies. Plesk is not responsible for support.<br />Please '
            . 'contact LiteSpeed at litespeedtech.com for all related '
            . 'questions and issues.<br /><br />LiteSpeed Web Server Extension '
            . 'for Plesk v'
            . IndexController::MODULE_VERSION
            . '</p></div>'
        ;

        $this->PageFooter();
    }

    /**
     *
     * @deprecated 1.4
     *
     * @param array $info
     */
    public function InstallLswsPrepare( array $info )
    {
        $this->PageHeader();

        $lswsHomeInputErr = $portOffsetErr = $adminEmailErr = $adminLoginErr
            = $adminPassErr = $adminPass1Err = null;

        if ( isset($info['error']['lsws_home_input']) ) {
            $lswsHomeInputErr = $info['error']['lsws_home_input'];
        }

        if ( isset($info['error']['port_offset']) ) {
            $portOffsetErr = $info['error']['port_offset'];
        }

        if ( isset($info['error']['admin_email']) ) {
            $adminEmailErr = $info['error']['admin_email'];
        }

        if ( isset($info['error']['admin_login']) ) {
            $adminLoginErr = $info['error']['admin_login'];
        }

        if ( isset($info['error']['admin_pass']) ) {
            $adminPassErr = $info['error']['admin_pass'];
        }

        if ( isset($info['error']['admin_pass1']) ) {
            $adminPass1Err = $info['error']['admin_pass1'];
        }

        $this->bufs[] =
            $this->screen_title(
                'Installing LiteSpeed Web Server',
                true,
                $this->icons['m_server_install']
            )
            . $this->show_choose_license($info)
            . $this->section_title('Installation Options')
            . $this->form_row(
                'Installation directory (define LSWS_HOME):',
                'text',
                $this->input_text(
                    'lsws_home_input',
                    $info['lsws_home_input'],
                    1
                ),
                $lswsHomeInputErr
            )
            . $this->form_row(
                'Port offset: ',
                'text',
                $this->input_text('port_offset', $info['port_offset']),
                $portOffsetErr,
                array(
                    'Setting a port offset allows you to run LiteSpeed on a '
                    . 'different port in parallel with Apache. The port '
                    . 'offset will be added to your Apache port number to '
                    . 'determine your LiteSpeed port.',
                    'It is recommended that you run LiteSpeed in parallel '
                    . 'first, so you can fully test it before switching to '
                    . 'LiteSpeed.'
                )
            )
            . $this->form_row(
                'Administrator email(s):',
                'text',
                $this->input_text('admin_email', $info['admin_email'], 2),
                $adminEmailErr,
                array(
                    '(Use commas to separate multiple email addresses.)',
                    'Email addresses specified will receive messages about '
                    . 'important events, such as server crashes or license '
                    . 'expirations.'
                )
            )
            . $this->section_title('WebAdmin Console Login')
            . $this->form_row(
                'User name:',
                'text',
                $this->input_text('admin_login', $info['admin_login']),
                $adminLoginErr
            )
            . $this->form_row(
                'Password:',
                'pass',
                $this->input_password('admin_pass', $info['admin_pass']),
                $adminPassErr
            )
            . $this->form_row(
                'Retype password:',
                'pass',
                $this->input_password('admin_pass1', $info['admin_pass1']),
                $adminPass1Err
            )
            . $this->button_panel_back_next('Cancel', 'Install')
        ;

        $this->PageFooter();
    }

    /**
     *
     * @deprecated 1.4
     *
     * @param array $info
     *
     * @throws PleskPluginException  Thrown indirectly by
     *     $this->show_running_status() call.
     */
    public function InstallLsws( array $info )
    {
        $this->PageHeader();

        $buf =
            $this->screen_title(
                'Install LiteSpeed Web Server',
                true,
                $this->icons['m_server_install']
            )
            . $this->show_running_status($info)
        ;

        if ( $info['return'] != 0 ) {
            $buf .= $this->error_panel_msg(
                'Error when installing LiteSpeed',
                $info['output']
            );
        }
        else {
            $buf .= (($info['spool_warning']) ? $this->getSpoolWarning() : '')
                . $this->info_panel_msg(
                    'LiteSpeed installed successfully',
                    $info['output']
                )
            ;
        }

        $this->bufs[] = $buf . $this->button_panel_back_next('OK');

        $this->PageFooter();
    }

    /**
     *
     * @deprecated 1.4
     *
     * @param array $info
     */
    public function DefineHome( array $info )
    {
        $this->PageHeader();

        $this->bufs[] =
            $this->screen_title(
                'Define LSWS_HOME Location for Existing LiteSpeed Installation',
                true,
                $this->icons['m_server_definehome']
            )
            . $this->info_msg(
                'If LiteSpeed is already installed on this server, please '
                . 'specify the LSWS_HOME location in order for this '
                . 'extension to work properly.'
            )
            . $this->section_title('Define $LSWS_HOME')
            . $this->form_row(
                '$LSWS_HOME location',
                'text',
                $this->input_text(
                    'lsws_home_input',
                    $info['lsws_home_input'],
                    1
                ),
                $info['error'],
                array(
                    'The LiteSpeed binary is located at '
                    . '$LSWS_HOME/bin/lshttpd.',
                    'Common locations for $LSWS_HOME include /usr/local/lsws '
                    . 'and /opt/lsws'
                )
            )
            . $this->button_panel_back_next('Cancel', 'Save')
        ;

        $this->PageFooter();
    }

}
